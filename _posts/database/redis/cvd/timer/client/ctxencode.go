package client

import (
	"context"

	"gitlab.xiaoduoai.com/golib/xd_sdk/metadata"
	"gitlab.xiaoduoai.com/golib/xd_sdk/octrace/propagation/textheader"
	"go.opencensus.io/trace"
)

func marshalSpanInfoToMD(ctx context.Context, md metadata.Metadata) context.Context {
	span := trace.FromContext(ctx)
	if span != nil {
		textheader.SpanContextToHeader(span.SpanContext(), md.GetD())
	}
	ctx = metadata.WithMetadata(ctx, md)
	return ctx
}

func unmarshalSpanInfoFromMD(ctx context.Context, md metadata.Metadata) (context.Context, func()) {
	sc, ok := textheader.SpanContextFromHeader(md.GetD())
	startOpts := trace.StartOptions{}
	if ok {
		var sampled trace.Sampler
		if sc.IsSampled() {
			sampled = trace.AlwaysSample()
		} else {
			sampled = trace.NeverSample()
		}

		startOpts = trace.StartOptions{
			Sampler: sampled,
		}
	}
	name := "timerwatch"

	var span *trace.Span
	if ok {
		ctx, span = trace.StartSpanWithRemoteParent(ctx, name, sc,
			trace.WithSampler(startOpts.Sampler),
			trace.WithSpanKind(trace.SpanKindServer))
	} else {
		ctx, span = trace.StartSpan(ctx, name,
			trace.WithSampler(startOpts.Sampler),
			trace.WithSpanKind(trace.SpanKindServer),
		)
	}
	return ctx, func() { span.End() }
}

func encodeCtx(ctx context.Context) string {
	md := metadata.FromContext(ctx)
	ctx = marshalSpanInfoToMD(ctx, md)

	_, mdstr := metadata.Encode(ctx)
	return mdstr
}

func decodeToCtx(mdStr string) (context.Context, func()) {
	ctx := metadata.FromString(context.Background(), mdStr)
	md := metadata.FromContext(ctx)

	ctx, end := unmarshalSpanInfoFromMD(ctx, md)
	return ctx, end
}
