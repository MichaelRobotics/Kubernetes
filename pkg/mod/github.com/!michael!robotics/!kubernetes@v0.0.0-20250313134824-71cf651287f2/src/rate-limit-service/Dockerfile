FROM python:3.11-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y protobuf-compiler && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY src/rate-limit-service/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

RUN mkdir -p /app/ratelimit/proto

RUN echo 'syntax = "proto3";\n\
package ratelimit;\n\
\n\
message RateLimitRequest {\n\
  // The domain is a unique name that represents a distributed ratelimit service.\n\
  string domain = 1;\n\
  // All requests must specify a descriptor.\n\
  repeated RateLimitDescriptor descriptors = 2;\n\
}\n\
\n\
message RateLimitDescriptor {\n\
  // A list of descriptor entries that make up the descriptor.\n\
  repeated RateLimitDescriptor.Entry entries = 1;\n\
  message Entry {\n\
    string key = 1;\n\
    string value = 2;\n\
  }\n\
}\n\
\n\
message RateLimitResponse {\n\
  // The overall response code.\n\
  enum Code {\n\
    UNKNOWN = 0;\n\
    OK = 1;\n\
    OVER_LIMIT = 2;\n\
  }\n\
  Code overall_code = 1;\n\
}\n\
\n\
service RateLimitService {\n\
  // Rate limits requests.\n\
  rpc ShouldRateLimit(RateLimitRequest) returns (RateLimitResponse) {}\n\
}\n\
' > /app/ratelimit/proto/ratelimit.proto

RUN cd /app && python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. ratelimit/proto/ratelimit.proto

RUN touch /app/ratelimit/__init__.py && \
    touch /app/ratelimit/proto/__init__.py

COPY src/rate-limit-service/app ./app
COPY src/rate-limit-service/config ./config

RUN adduser --disabled-password --gecos "" appuser
USER appuser

ENV PYTHONUNBUFFERED=1

EXPOSE 8081 8082 9090

ENTRYPOINT [ "opentelemetry-instrument", "python", "-m", "app.main" ] 