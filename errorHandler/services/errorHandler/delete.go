package errorHandler

import (
	"context"


	erHan "errorHandler/gunk/v1/errorHandler"
)

func (s *Svc) DeleteErrorHandler(ctx context.Context, req *erHan.DeleteErrorHandlerRequest) (*erHan.DeleteErrorHandlerResponse, error) {
	
	if err := s.erHan.DeleteErrorHandler(ctx, req.ID, req.DeleteByEnvType); err != nil {
		
		return nil, err
	}

	return &erHan.DeleteErrorHandlerResponse{}, nil
}
