package main

const postTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .Entry.Title}}{{.Entry.Title}} - {{end}}{{.SiteTitle}}</title>
    <style>
{{.ThemeCSS}}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="header-content">
                <div>
                    <h1><a href="/" style="color: inherit; text-decoration: none;">{{.SiteTitle}}</a></h1>
                    <div class="subtitle">{{.SiteSubtitle}}</div>
                </div>
            </div>
        </div>

        <div class="feed">
            <div class="entry">
                <div class="entry-header">
                    <div class="avatar">
                        {{if and (eq .AvatarPreference "avatar") .AvatarPath}}
                            <img src="/uploads/{{.AvatarPath}}" alt="Avatar" style="width: 100%; height: 100%; border-radius: 50%; object-fit: cover;">
                        {{else}}
                            <div class="avatar-inner">{{.UserInitial}}</div>
                        {{end}}
                    </div>
                </div>
                {{if .Entry.HasPhoto}}
                <div class="entry-photo-container">
                    <img src="{{.Entry.Photo}}" alt="Entry photo" class="entry-photo">
                </div>
                {{end}}
                {{if .Entry.HasAudio}}
                <div class="entry-photo-container">
                    {{if .Entry.HasThumbnail}}
                    <img src="{{.Entry.Thumbnail}}" alt="Audio cover" class="entry-photo" style="object-fit: cover;">
                    {{else}}
                    <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center;">
                        <svg width="80" height="80" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="opacity: 0.7;">
                            <path d="M9 18V5l12-2v13"></path>
                            <circle cx="6" cy="18" r="3"></circle>
                            <circle cx="18" cy="16" r="3"></circle>
                        </svg>
                    </div>
                    {{end}}
                    <div style="position: absolute; bottom: 0; left: 0; right: 0; padding: 16px; background: linear-gradient(transparent, rgba(0,0,0,0.8));">
                        <audio controls style="width: 100%; filter: invert(1) hue-rotate(180deg);">
                            <source src="{{.Entry.Photo}}">
                            Your browser does not support the audio element.
                        </audio>
                    </div>
                </div>
                {{end}}
                {{if .Entry.HasVideo}}
                <div class="entry-photo-container">
                    {{if .Entry.HasThumbnail}}
                    <div class="video-thumbnail-wrapper" id="videoThumbnail" onclick="document.getElementById('videoThumbnail').style.display='none'; document.getElementById('videoPlayer').style.display='block'; document.getElementById('videoPlayer').querySelector('video').play();">
                        <img src="{{.Entry.Thumbnail}}" alt="Video thumbnail" class="entry-photo" style="object-fit: cover; cursor: pointer;">
                        <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); background: rgba(0,0,0,0.6); border-radius: 50%; width: 80px; height: 80px; display: flex; align-items: center; justify-content: center; cursor: pointer;">
                            <svg width="40" height="40" viewBox="0 0 24 24" fill="white">
                                <polygon points="5 3 19 12 5 21 5 3"></polygon>
                            </svg>
                        </div>
                    </div>
                    <div id="videoPlayer" style="display: none;">
                        <video controls class="entry-photo" style="object-fit: contain;">
                            <source src="{{.Entry.Photo}}" type="video/mp4">
                            Your browser does not support the video element.
                        </video>
                    </div>
                    {{else}}
                    <video controls class="entry-photo" style="object-fit: contain;">
                        <source src="{{.Entry.Photo}}" type="video/mp4">
                        Your browser does not support the video element.
                    </video>
                    {{end}}
                </div>
                {{end}}
                <div class="entry-content">
                    {{if .Entry.IsTruncated}}
                        <span id="truncated-content">{{.Entry.Content}}</span>
                        <span id="full-content" style="display: none;">{{.Entry.FullContent}}</span>
                        <a href="#" id="read-more-link" style="color: #0095f6; text-decoration: none; font-weight: 500; margin-left: 4px;">Read more</a>
                    {{else}}
                        {{.Entry.Content}}
                    {{end}}
                </div>
                <div class="entry-timestamp">{{.Entry.TimeAgo}}</div>
            </div>

            <div style="text-align: center; margin-top: 32px;">
                <a href="/" style="color: #0095f6; text-decoration: none; font-size: 14px; font-weight: 600;">← Back to all posts</a>
            </div>
        </div>
    </div>

    <script>
        document.addEventListener('DOMContentLoaded', function() {
            const readMoreLink = document.getElementById('read-more-link');
            if (readMoreLink) {
                readMoreLink.addEventListener('click', function(e) {
                    e.preventDefault();
                    const truncatedContent = document.getElementById('truncated-content');
                    const fullContent = document.getElementById('full-content');

                    if (fullContent.style.display === 'none') {
                        truncatedContent.style.display = 'none';
                        fullContent.style.display = 'inline';
                        this.textContent = 'Show less';
                    } else {
                        truncatedContent.style.display = 'inline';
                        fullContent.style.display = 'none';
                        this.textContent = 'Read more';
                    }
                });
            }
        });
    </script>
</body>
</html>`
