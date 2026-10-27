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
	mux.HandleFunc("GET /api/projects", h.getProjects)
	mux.HandleFunc("GET /api/projects/featured", h.getFeaturedProjects)
	mux.HandleFunc("GET /api/projects/{id}", h.getProject)
	mux.HandleFunc("GET /api/experiences", h.getExperiences)
	mux.HandleFunc("GET /api/skill", h.getSkills)
	return mux, nil
}

func (h *getHandlers) getProjects(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, r, "/api/projects")
}

func (h *getHandlers) getFeaturedProjects(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, r, "/api/projects/featured")
}

func (h *getHandlers) getProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	projectID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || projectID < 1 {
		http.Error(w, "project id must be a positive integer", http.StatusBadRequest)
		return
	}
	h.proxy(w, r, "/api/projects/"+strconv.FormatInt(projectID, 10))
}

func (h *getHandlers) getExperiences(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, r, "/api/experiences")
}

func (h *getHandlers) getSkills(w http.ResponseWriter, r *http.Request) {
	h.proxy(w, r, "/api/skill")
}

func (h *getHandlers) proxy(w http.ResponseWriter, r *http.Request, endpoint string) {
	upstreamURL := *h.upstream
	upstreamURL.Path = strings.TrimRight(upstreamURL.Path, "/") + endpoint
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
}
