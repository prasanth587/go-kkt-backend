# Fix for Image Storage on Railway

## Problem
Railway's filesystem is **ephemeral** - all files in `/app/uploads/` are deleted when you redeploy. The database still has the image paths, but the actual files are gone.

## Solutions

### Option 1: Railway Volumes (Recommended - Easiest)

1. **In Railway Dashboard:**
   - Go to your service
   - Click on "Volumes" tab
   - Click "New Volume"
   - Name it: `uploads-volume`
   - Set mount path: `/app/uploads`
   - Create the volume

2. **Update Environment Variables:**
   - Make sure `BASE_DIRECTORY=/app/uploads` is set
   - Make sure `IMAGE_DIRECTORY=/t_hub_document` is set

3. **Redeploy:**
   - The volume will persist across redeploys
   - Images will survive redeployments

### Option 2: Cloud Storage (AWS S3, Google Cloud Storage, etc.)

This requires code changes to upload files to cloud storage instead of local filesystem.

**Benefits:**
- Files are stored in the cloud (more reliable)
- Can be accessed from anywhere
- Better for scaling

**Implementation:**
- Use AWS S3 SDK or Google Cloud Storage SDK
- Upload files to cloud storage bucket
- Store the cloud URL in database instead of local path
- Serve images directly from cloud storage or through a CDN

### Option 3: Accept the Limitation (Not Recommended)

- Users need to re-upload images after each redeploy
- Only suitable for development/testing

## Current Status

The error message has been improved to inform users: 
> "Image file not found. This may happen if the server was redeployed. Please re-upload the image."

## Recommendation

**Use Railway Volumes (Option 1)** - It's the quickest fix with minimal code changes.

