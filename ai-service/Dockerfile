FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application
COPY . .

# Create data directory
RUN mkdir -p /app/data/chroma_db
RUN mkdir -p /app/logs

# Expose ports
EXPOSE 8090 50070

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD python -c "import requests; requests.get('http://localhost:8090/v2/ai/health')"

# Run application
CMD ["python", "main.py"]

