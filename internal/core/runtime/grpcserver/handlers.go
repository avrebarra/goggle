package grpcserver

import (
	"context"
	"errors"
	"time"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/internal/module/servicetoggle/domaintoggle"
	pb "github.com/avrebarra/goggle/proto/generated"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/guregu/null/v5"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	ConfigRuntime
	pb.UnimplementedGoggleServiceServer
}

// Ping implements the health check endpoint
func (h *Handler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Version:   h.Version,
		StartedAt: timestamppb.New(h.StartedAt),
		Uptime:    time.Since(h.StartedAt).Round(time.Second).String(),
	}, nil
}

// ListToggles lists all toggles with pagination and filtering
func (h *Handler) ListToggles(ctx context.Context, req *pb.ListTogglesRequest) (*pb.ListTogglesResponse, error) {
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, pkgerrors.Wrap(err, "bad request").Error())
	}

	params := servicetoggle.ParamsDoListToggles{
		Offset:    int(req.Offset),
		Limit:     int(req.Limit),
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.FilterAccessed != nil {
		params.FilterAccessed = null.BoolFrom(req.FilterAccessed.Value)
	}

	data, total, err := h.ToggleService.DoListToggles(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, pkgerrors.Wrap(err, "service failure").Error())
	}

	items := make([]*pb.ToggleWithDetail, 0, len(data))
	for _, d := range data {
		item := &pb.ToggleWithDetail{
			Id:               d.ID,
			Status:           d.Status,
			UpdatedAt:        timestamppb.New(d.UpdatedAt),
			AccessFreqWeekly: int32(d.AccessFreqWeekly),
		}
		if !d.LastAccessedAt.IsZero() {
			item.LastAccessedAt = timestamppb.New(d.LastAccessedAt)
		}
		items = append(items, item)
	}

	return &pb.ListTogglesResponse{
		Items: items,
		Total: total,
	}, nil
}

// ListStrayToggles lists stray (unused/orphaned) toggles
func (h *Handler) ListStrayToggles(ctx context.Context, req *pb.ListStrayTogglesRequest) (*pb.ListStrayTogglesResponse, error) {
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, pkgerrors.Wrap(err, "bad request").Error())
	}

	params := servicetoggle.ParamsDoListStrayToggles{
		Offset:    int(req.Offset),
		Limit:     int(req.Limit),
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	data, total, err := h.ToggleService.DoListStrayToggles(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, pkgerrors.Wrap(err, "service failure").Error())
	}

	items := make([]*pb.ToggleStatLog, 0, len(data))
	for _, d := range data {
		item := &pb.ToggleStatLog{
			Id:               d.ID,
			AccessFreqWeekly: int32(d.AccessFreqWeekly),
		}
		if !d.LastAccessedAt.IsZero() {
			item.LastAccessedAt = timestamppb.New(d.LastAccessedAt)
		}
		items = append(items, item)
	}

	return &pb.ListStrayTogglesResponse{
		Items: items,
		Total: total,
	}, nil
}

// GetToggle retrieves a specific toggle by ID
func (h *Handler) GetToggle(ctx context.Context, req *pb.GetToggleRequest) (*pb.ToggleWithDetail, error) {
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, pkgerrors.Wrap(err, "bad request").Error())
	}

	data, err := h.ToggleService.DoGetToggle(ctx, req.Id)
	if errors.Is(err, servicetoggle.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "toggle not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, pkgerrors.Wrap(err, "service failure").Error())
	}

	result := &pb.ToggleWithDetail{
		Id:               data.ID,
		Status:           data.Status,
		UpdatedAt:        timestamppb.New(data.UpdatedAt),
		AccessFreqWeekly: int32(data.AccessFreqWeekly),
	}
	if !data.LastAccessedAt.IsZero() {
		result.LastAccessedAt = timestamppb.New(data.LastAccessedAt)
	}

	return result, nil
}

// CreateToggle creates a new toggle
func (h *Handler) CreateToggle(ctx context.Context, req *pb.CreateToggleRequest) (*pb.Toggle, error) {
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, pkgerrors.Wrap(err, "bad request").Error())
	}

	toggle := domaintoggle.Toggle{
		ID:        req.Id,
		Status:    req.Status,
		UpdatedAt: req.UpdatedAt.AsTime(),
	}

	data, err := h.ToggleService.DoCreateToggle(ctx, toggle)
	if err != nil {
		return nil, status.Error(codes.Internal, pkgerrors.Wrap(err, "service failure").Error())
	}

	return &pb.Toggle{
		Id:        data.ID,
		Status:    data.Status,
		UpdatedAt: timestamppb.New(data.UpdatedAt),
	}, nil
}

// UpdateToggle updates an existing toggle
func (h *Handler) UpdateToggle(ctx context.Context, req *pb.UpdateToggleRequest) (*pb.Toggle, error) {
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, pkgerrors.Wrap(err, "bad request").Error())
	}

	toggle := domaintoggle.Toggle{
		ID:        req.Data.Id,
		Status:    req.Data.Status,
		UpdatedAt: req.Data.UpdatedAt.AsTime(),
	}

	data, err := h.ToggleService.DoUpdateToggle(ctx, req.Id, toggle)
	if err != nil {
		return nil, status.Error(codes.Internal, pkgerrors.Wrap(err, "service failure").Error())
	}

	return &pb.Toggle{
		Id:        data.ID,
		Status:    data.Status,
		UpdatedAt: timestamppb.New(data.UpdatedAt),
	}, nil
}

// RemoveToggle deletes a toggle
func (h *Handler) RemoveToggle(ctx context.Context, req *pb.RemoveToggleRequest) (*pb.Toggle, error) {
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, pkgerrors.Wrap(err, "bad request").Error())
	}

	data, err := h.ToggleService.DoRemoveToggle(ctx, req.Id)
	if errors.Is(err, servicetoggle.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "toggle not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, pkgerrors.Wrap(err, "service failure").Error())
	}

	return &pb.Toggle{
		Id:        data.ID,
		Status:    data.Status,
		UpdatedAt: timestamppb.New(data.UpdatedAt),
	}, nil
}

// StatToggle retrieves toggle statistics
func (h *Handler) StatToggle(ctx context.Context, req *pb.StatToggleRequest) (*pb.ToggleCompact, error) {
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, pkgerrors.Wrap(err, "bad request").Error())
	}

	data, err := h.ToggleService.DoStatToggle(ctx, req.Id)
	if errors.Is(err, servicetoggle.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "toggle not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, pkgerrors.Wrap(err, "service failure").Error())
	}

	return &pb.ToggleCompact{
		Id:     data.ID,
		Status: data.Status,
	}, nil
}
