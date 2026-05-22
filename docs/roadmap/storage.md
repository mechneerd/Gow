> ✅ Implemented · 🚧 In Progress · 📋 Planned

# File Storage

> **Status**: 🚧 In Progress (Local disk implemented, S3 stubbed)


GoW provides a powerful filesystem abstraction using a multi-disk `Storage` manager. This allows you to easily swap between local filesystems and cloud storage like Amazon S3 without changing your application code.

## Configuration

Your filesystem configuration is defined in `config/filesystems.go`. Within this file, you may configure all of your "disks". Each disk represents a particular storage driver and storage location.

## Basic Usage

The `Storage` manager provides an expressive interface to interact with files.

```go
// Store a file on the default disk
err := storage.Put("avatars/user.jpg", fileReader)

// Retrieve a file
reader, err := storage.Get("avatars/user.jpg")

// Check if a file exists
exists := storage.Exists("avatars/user.jpg")

// Delete a file
err := storage.Delete("avatars/user.jpg")
```

## Using Specific Disks

If your application interacts with multiple disks, you may use the `Disk` method to work with files on a particular disk.

```go
storage.Disk("s3").Put("avatars/user.jpg", fileReader)
```

## File URLs

You may use the `URL` method to get the URL for a given file. If you are using the `local` driver, this will typically prepend the `/storage` path to the given path. If you are using the `s3` driver, the fully qualified AWS S3 URL will be returned.

```go
url := storage.URL("avatars/user.jpg")
```
