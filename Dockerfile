FROM golang:latest

# Set destination for COPY
WORKDIR /app

# install python/pip
RUN apt-get update && apt-get install -y \ 
    python3 \
    python3-pip \ 
    python3-venv 
# Create/activate python VE to go from apt-get to pip package management to avoid:
# error - externally-managed-environment 
RUN python3 -m venv /opt/venv 
ENV PATH="/opt/venv/bin:$PATH"
# Download python packages
COPY /transcribe/requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY /handler/ /app/handler/
COPY /render/ /app/render/
COPY /public/ /app/public/
COPY /transcribe/ /app/transcribe/
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /prattl

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8081

# Run
CMD ["/prattl"]
