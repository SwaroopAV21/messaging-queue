package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/SwaroopAV21/messaging-queue/models"
	"github.com/SwaroopAV21/messaging-queue/repository"
)


type getHandlers struct {
	upstream *url.URL
	client   *http.Client
}

func NewGetHandler(upstreamBaseURL string) (http.Handler, error) {
	upstream, err := url.Parse(strings.TrimRight(upstreamBaseURL, "/"))
	if err != nil {
		return nil, fmt.Errorf("parse upstream base URL: %w", err)
	}
	if (upstream.Scheme != "http" && upstream.Scheme != "https") || upstream.Host == "" ||
		upstream.User != nil || upstream.RawQuery != "" || upstream.Fragment != "" {
		return nil, fmt.Errorf("upstream base URL must be an HTTP(S) URL without credentials, query, or fragment")
	}

	h := &getHandlers{
		upstream: upstream,
		client:   &http.Client{Timeout: 15 * time.Second},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/middleware/projects", h.getProjects)
	mux.HandleFunc("GET /api/middleware/projects/featured", h.getFeaturedProjects)
	mux.HandleFunc("GET /api/middleware/projects/{id}", h.getProject)
	mux.HandleFunc("GET /api/middleware/experiences", h.getExperiences)
	mux.HandleFunc("GET /api/middleware/skill", h.getSkills)
	mux.HandleFunc("GET /api/middleware/responses", getResponses)
	mux.HandleFunc("GET /api/middleware/requests", getRequests)
	return mux, nil
}

func (h *getHandlers) getProjects(w http.ResponseWriter, r *http.Request) {
	requestProjects := models.Request {
		Api: "/api/projects",
		Method: r.Method,
	}
	h.proxy(w, r, requestProjects)
}

func (h *getHandlers) getFeaturedProjects(w http.ResponseWriter, r *http.Request) {
	requestFeaturedProjects := models.Request {
		Api: "/api/projects/featured",
		Method: r.Method,
	}
	h.proxy(w, r, requestFeaturedProjects)
}

func (h *getHandlers) getProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	projectID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || projectID < 1 {
		http.Error(w, "project id must be a positive integer", http.StatusBadRequest)
		return
	}

	requestProject := models.Request {
		Api: "/api/projects/" + strconv.FormatInt(projectID, 10),
		Method: r.Method,
	}

	h.proxy(w, r, requestProject)
}

func (h *getHandlers) getExperiences(w http.ResponseWriter, r *http.Request) {
	
	requestExperiences := models.Request {
		Api: "/api/experiences",
		Method: r.Method,
	}
	h.proxy(w, r, requestExperiences)
}

func (h *getHandlers) getSkills(w http.ResponseWriter, r *http.Request) {
	
	requestSkills := models.Request {
		Api: "/api/skills/grouped",
		Method: r.Method,
	}
	h.proxy(w, r, requestSkills)
}

func (h *getHandlers) proxy(w http.ResponseWriter, r *http.Request, requestObj models.Request) {
	upstreamURL := *h.upstream
	upstreamURL.Path = strings.TrimRight(upstreamURL.Path, "/") + requestObj.Api
	upstreamURL.RawQuery = r.URL.RawQuery

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, upstreamURL.String(), nil)
	if err != nil {
		http.Error(w, "could not create upstream request", http.StatusInternalServerError)
		return
	}
	if accept := r.Header.Get("Accept"); accept != "" {
		req.Header.Set("Accept", accept)
	}

	resp, err := h.client.Do(req)


	if err != nil {
		http.Error(w, "upstream service unavailable", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	if cacheControl := resp.Header.Get("Cache-Control"); cacheControl != "" {
		w.Header().Set("Cache-Control", cacheControl)
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("copy upstream response body: %v", err)
	}

	responseObj := models.Response {
		RequestID: requestObj.RequestID,
		Api: requestObj.Api,         
		ReturnCode: resp.StatusCode,
		RequestStatus: resp.Status,
	}

	if err := repository.SaveToDB(r.Context(), responseObj); err != nil {
    	log.Printf("save response: %v", err)
    	http.Error(w, "could not save response", http.StatusInternalServerError)
    	return
	}

}
