/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/golang/glog"
)

// only allow pods to pull images from specific registry.
func admit(data []byte) *AdmissionReviewStatus {
	ar := AdmissionReview{}
	if err := json.Unmarshal(data, &ar); err != nil {
		glog.Error(err)
		return nil
	}
	// The externalAdmissionHookConfiguration registered via selfRegistration
	// asks the kube-apiserver only sends admission request regarding pods.
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Spec.Resource != podResource {
		glog.Errorf("expect resource to be %s", podResource)
		return nil
	}

	raw := ar.Spec.Object.Raw
	pod := v1.Pod{}
	if err := json.Unmarshal(raw, &pod); err != nil {
		glog.Error(err)
		return nil
	}
	reviewStatus := AdmissionReviewStatus{}
	for _, container := range pod.Spec.Containers {
		// gcr.io is just an example.
		if !strings.Contains(container.Image, "gcr.io") {
			reviewStatus.Allowed = false
			reviewStatus.Result = &metav1.Status{
				Reason: "can only pull image from grc.io",
			}
			return &reviewStatus
		}
	}
	reviewStatus.Allowed = true
	return &reviewStatus
}

func serve(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	reviewStatus := admit(body)
	ar := AdmissionReview{
		Status: *reviewStatus,
	}

	resp, err := json.Marshal(ar)
	if err != nil {
		glog.Error(err)
	}
	if _, err := w.Write(resp); err != nil {
		glog.Error(err)
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", serve)
	clientset := getClient()
	server := &http.Server{
		Addr:      ":8000",
		TLSConfig: configTLS(clientset),
	}
	go selfRegistration(clientset, caCert)
	server.ListenAndServeTLS("", "")
}
