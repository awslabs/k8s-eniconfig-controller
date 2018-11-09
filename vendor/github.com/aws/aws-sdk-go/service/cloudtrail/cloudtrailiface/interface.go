// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package cloudtrailiface provides an interface to enable mocking the AWS CloudTrail service client
// for testing your code.
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters.
package cloudtrailiface

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
)

// CloudTrailAPI provides an interface to enable mocking the
// cloudtrail.CloudTrail service client's API operation,
// paginators, and waiters. This make unit testing your code that calls out
// to the SDK's service client's calls easier.
//
// The best way to use this interface is so the SDK's service client's calls
// can be stubbed out for unit testing your code with the SDK without needing
// to inject custom request handlers into the SDK's request pipeline.
//
//    // myFunc uses an SDK service client to make a request to
//    // AWS CloudTrail.
//    func myFunc(svc cloudtrailiface.CloudTrailAPI) bool {
//        // Make svc.AddTags request
//    }
//
//    func main() {
//        sess := session.New()
//        svc := cloudtrail.New(sess)
//
//        myFunc(svc)
//    }
//
// In your _test.go file:
//
//    // Define a mock struct to be used in your unit tests of myFunc.
//    type mockCloudTrailClient struct {
//        cloudtrailiface.CloudTrailAPI
//    }
//    func (m *mockCloudTrailClient) AddTags(input *cloudtrail.AddTagsInput) (*cloudtrail.AddTagsOutput, error) {
//        // mock response/functionality
//    }
//
//    func TestMyFunc(t *testing.T) {
//        // Setup Test
//        mockSvc := &mockCloudTrailClient{}
//
//        myfunc(mockSvc)
//
//        // Verify myFunc's functionality
//    }
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters. Its suggested to use the pattern above for testing, or using
// tooling to generate mocks to satisfy the interfaces.
type CloudTrailAPI interface {
	AddTags(*cloudtrail.AddTagsInput) (*cloudtrail.AddTagsOutput, error)
	AddTagsWithContext(aws.Context, *cloudtrail.AddTagsInput, ...request.Option) (*cloudtrail.AddTagsOutput, error)
	AddTagsRequest(*cloudtrail.AddTagsInput) (*request.Request, *cloudtrail.AddTagsOutput)

	CreateTrail(*cloudtrail.CreateTrailInput) (*cloudtrail.CreateTrailOutput, error)
	CreateTrailWithContext(aws.Context, *cloudtrail.CreateTrailInput, ...request.Option) (*cloudtrail.CreateTrailOutput, error)
	CreateTrailRequest(*cloudtrail.CreateTrailInput) (*request.Request, *cloudtrail.CreateTrailOutput)

	DeleteTrail(*cloudtrail.DeleteTrailInput) (*cloudtrail.DeleteTrailOutput, error)
	DeleteTrailWithContext(aws.Context, *cloudtrail.DeleteTrailInput, ...request.Option) (*cloudtrail.DeleteTrailOutput, error)
	DeleteTrailRequest(*cloudtrail.DeleteTrailInput) (*request.Request, *cloudtrail.DeleteTrailOutput)

	DescribeTrails(*cloudtrail.DescribeTrailsInput) (*cloudtrail.DescribeTrailsOutput, error)
	DescribeTrailsWithContext(aws.Context, *cloudtrail.DescribeTrailsInput, ...request.Option) (*cloudtrail.DescribeTrailsOutput, error)
	DescribeTrailsRequest(*cloudtrail.DescribeTrailsInput) (*request.Request, *cloudtrail.DescribeTrailsOutput)

	GetEventSelectors(*cloudtrail.GetEventSelectorsInput) (*cloudtrail.GetEventSelectorsOutput, error)
	GetEventSelectorsWithContext(aws.Context, *cloudtrail.GetEventSelectorsInput, ...request.Option) (*cloudtrail.GetEventSelectorsOutput, error)
	GetEventSelectorsRequest(*cloudtrail.GetEventSelectorsInput) (*request.Request, *cloudtrail.GetEventSelectorsOutput)

	GetTrailStatus(*cloudtrail.GetTrailStatusInput) (*cloudtrail.GetTrailStatusOutput, error)
	GetTrailStatusWithContext(aws.Context, *cloudtrail.GetTrailStatusInput, ...request.Option) (*cloudtrail.GetTrailStatusOutput, error)
	GetTrailStatusRequest(*cloudtrail.GetTrailStatusInput) (*request.Request, *cloudtrail.GetTrailStatusOutput)

	ListPublicKeys(*cloudtrail.ListPublicKeysInput) (*cloudtrail.ListPublicKeysOutput, error)
	ListPublicKeysWithContext(aws.Context, *cloudtrail.ListPublicKeysInput, ...request.Option) (*cloudtrail.ListPublicKeysOutput, error)
	ListPublicKeysRequest(*cloudtrail.ListPublicKeysInput) (*request.Request, *cloudtrail.ListPublicKeysOutput)

	ListTags(*cloudtrail.ListTagsInput) (*cloudtrail.ListTagsOutput, error)
	ListTagsWithContext(aws.Context, *cloudtrail.ListTagsInput, ...request.Option) (*cloudtrail.ListTagsOutput, error)
	ListTagsRequest(*cloudtrail.ListTagsInput) (*request.Request, *cloudtrail.ListTagsOutput)

	LookupEvents(*cloudtrail.LookupEventsInput) (*cloudtrail.LookupEventsOutput, error)
	LookupEventsWithContext(aws.Context, *cloudtrail.LookupEventsInput, ...request.Option) (*cloudtrail.LookupEventsOutput, error)
	LookupEventsRequest(*cloudtrail.LookupEventsInput) (*request.Request, *cloudtrail.LookupEventsOutput)

	LookupEventsPages(*cloudtrail.LookupEventsInput, func(*cloudtrail.LookupEventsOutput, bool) bool) error
	LookupEventsPagesWithContext(aws.Context, *cloudtrail.LookupEventsInput, func(*cloudtrail.LookupEventsOutput, bool) bool, ...request.Option) error

	PutEventSelectors(*cloudtrail.PutEventSelectorsInput) (*cloudtrail.PutEventSelectorsOutput, error)
	PutEventSelectorsWithContext(aws.Context, *cloudtrail.PutEventSelectorsInput, ...request.Option) (*cloudtrail.PutEventSelectorsOutput, error)
	PutEventSelectorsRequest(*cloudtrail.PutEventSelectorsInput) (*request.Request, *cloudtrail.PutEventSelectorsOutput)

	RemoveTags(*cloudtrail.RemoveTagsInput) (*cloudtrail.RemoveTagsOutput, error)
	RemoveTagsWithContext(aws.Context, *cloudtrail.RemoveTagsInput, ...request.Option) (*cloudtrail.RemoveTagsOutput, error)
	RemoveTagsRequest(*cloudtrail.RemoveTagsInput) (*request.Request, *cloudtrail.RemoveTagsOutput)

	StartLogging(*cloudtrail.StartLoggingInput) (*cloudtrail.StartLoggingOutput, error)
	StartLoggingWithContext(aws.Context, *cloudtrail.StartLoggingInput, ...request.Option) (*cloudtrail.StartLoggingOutput, error)
	StartLoggingRequest(*cloudtrail.StartLoggingInput) (*request.Request, *cloudtrail.StartLoggingOutput)

	StopLogging(*cloudtrail.StopLoggingInput) (*cloudtrail.StopLoggingOutput, error)
	StopLoggingWithContext(aws.Context, *cloudtrail.StopLoggingInput, ...request.Option) (*cloudtrail.StopLoggingOutput, error)
	StopLoggingRequest(*cloudtrail.StopLoggingInput) (*request.Request, *cloudtrail.StopLoggingOutput)

	UpdateTrail(*cloudtrail.UpdateTrailInput) (*cloudtrail.UpdateTrailOutput, error)
	UpdateTrailWithContext(aws.Context, *cloudtrail.UpdateTrailInput, ...request.Option) (*cloudtrail.UpdateTrailOutput, error)
	UpdateTrailRequest(*cloudtrail.UpdateTrailInput) (*request.Request, *cloudtrail.UpdateTrailOutput)
}

var _ CloudTrailAPI = (*cloudtrail.CloudTrail)(nil)
