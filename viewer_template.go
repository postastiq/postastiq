package main

const defaultThemeCSS = `
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: #fafafa;
            color: #262626;
            line-height: 1.6;
        }

        .container {
            max-width: 614px;
            margin: 0 auto;
            background-color: #fafafa;
            min-height: 100vh;
        }

        .header {
            background: #ffffff;
            padding: 20px 32px;
            color: #262626;
            border-bottom: 1px solid #dbdbdb;
            position: sticky;
            top: 0;
            z-index: 100;
        }

        .header-content {
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .header h1 {
            font-size: 28px;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 12px;
            font-family: 'Segoe UI', Roboto, sans-serif;
        }

        .subtitle {
            font-size: 14px;
            color: #8e8e8e;
            margin-top: 4px;
            font-weight: 400;
        }

        .stats {
            padding: 16px 0;
            background-color: #fafafa;
            display: none;
        }

        .stat {
            text-align: center;
            padding: 16px;
            background-color: #ffffff;
            border: 1px solid #dbdbdb;
            border-radius: 8px;
            transition: transform 0.2s, box-shadow 0.2s;
        }

        .stat:hover {
            transform: translateY(-2px);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }

        .stat-value {
            font-size: 32px;
            font-weight: 700;
            color: #262626;
            margin-bottom: 4px;
        }

        .stat-label {
            font-size: 13px;
            color: #8e8e8e;
            text-transform: uppercase;
            font-weight: 500;
            letter-spacing: 0.5px;
        }

        .feed {
            padding: 24px 0 80px;
        }

        .entry {
            padding: 0;
            margin-bottom: 24px;
            background-color: #ffffff;
            border: 1px solid #dbdbdb;
            border-radius: 8px;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            display: flex;
            flex-direction: column;
        }

        .entry:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }

        .entry-header {
            display: flex;
            align-items: center;
            padding: 14px 16px;
            border-bottom: 1px solid #efefef;
        }

        .avatar {
            width: 56px;
            height: 56px;
            border-radius: 50%;
            background: #000000;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 600;
            font-size: 20px;
            color: #ffffff;
            margin-right: 16px;
            flex-shrink: 0;
            padding: 3px;
        }

        .avatar-inner {
            width: 100%;
            height: 100%;
            background: #ffffff;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #262626;
        }

        .entry-info {
            flex: 1;
        }

        .username {
            display: none;
        }

        .timestamp {
            color: #8e8e8e;
            font-size: 12px;
            font-weight: 400;
            margin-top: 4px;
            display: block;
        }

        .entry-photo-container {
            width: 100%;
            max-width: 1024px;
            height: 0;
            padding-bottom: 100%;
            position: relative;
            overflow: hidden;
            background-color: #000000;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .entry-photo {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            max-width: 100%;
            max-height: 100%;
            width: auto;
            height: auto;
            object-fit: contain;
            display: block;
        }

        .entry-actions {
            display: none;
        }

        .action-btn {
            background: none;
            border: none;
            cursor: pointer;
            font-size: 24px;
            padding: 8px;
            line-height: 1;
            color: #262626;
            transition: opacity 0.2s;
        }

        .action-btn:hover {
            opacity: 0.6;
        }

        .entry-content {
            font-size: 15px;
            color: #262626;
            line-height: 20px;
            padding: 16px 16px 4px;
            word-wrap: break-word;
            overflow-wrap: break-word;
        }

        .entry-content a {
            color: #0095f6;
            text-decoration: none;
            font-weight: 500;
        }

        .entry-content a:hover {
            text-decoration: underline;
        }

        .entry-content-username {
            display: none;
        }

        .entry-timestamp {
            padding: 4px 16px 16px;
            color: #8e8e8e;
            font-size: 10px;
            text-transform: uppercase;
            letter-spacing: 0.2px;
        }

        .entry-footer {
            display: none;
        }

        .entry-id {
            color: #8e8e8e;
            font-size: 12px;
            font-weight: 500;
        }

        .entry-badge {
            background-color: #0095f6;
            color: #ffffff;
            padding: 4px 12px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .empty-state {
            padding: 80px 32px;
            text-align: center;
        }

        .empty-state-icon {
            width: 120px;
            height: 120px;
            margin: 0 auto 24px;
            background: #ffffff;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 64px;
            border: 2px solid #dbdbdb;
        }

        .empty-state-title {
            font-size: 28px;
            font-weight: 600;
            margin-bottom: 12px;
            color: #262626;
        }

        .empty-state-text {
            font-size: 16px;
            color: #8e8e8e;
            max-width: 400px;
            margin: 0 auto;
        }

        .loading-indicator {
            text-align: center;
            padding: 24px;
            color: #8e8e8e;
            font-size: 14px;
            display: none;
        }

        .loading-indicator.show {
            display: block;
        }

        .loading-spinner {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 2px solid #dbdbdb;
            border-top-color: #262626;
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
            margin-right: 8px;
            vertical-align: middle;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        .end-message {
            text-align: center;
            padding: 24px;
            color: #8e8e8e;
            font-size: 14px;
            display: none;
        }

        .end-message.show {
            display: block;
        }

        @media (max-width: 768px) {
            .container {
                background-color: #ffffff;
            }

            .header {
                padding: 16px 20px;
            }

            .header h1 {
                font-size: 24px;
            }

            .stats {
                padding: 12px 16px;
                gap: 12px;
            }

            .stat {
                padding: 12px;
            }

            .stat-value {
                font-size: 24px;
            }

            .feed {
                padding: 0 0 80px;
            }

            .entry {
                margin-bottom: 12px;
                border-radius: 0;
                border-left: none;
                border-right: none;
            }
        }
`

const darkThemeCSS = `
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: #0a0a0a;
            color: #e4e6eb;
            line-height: 1.6;
        }

        .container {
            max-width: 614px;
            margin: 0 auto;
            background-color: #0a0a0a;
            min-height: 100vh;
        }

        .header {
            background: #1c1e21;
            padding: 20px 32px;
            color: #e4e6eb;
            border-bottom: 1px solid #2f3336;
            position: sticky;
            top: 0;
            z-index: 100;
        }

        .header-content {
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .header h1 {
            font-size: 28px;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 12px;
            font-family: 'Segoe UI', Roboto, sans-serif;
        }

        .subtitle {
            font-size: 14px;
            color: #8b949e;
            margin-top: 4px;
            font-weight: 400;
        }

        .stats {
            padding: 16px 0;
            background-color: #0a0a0a;
            display: none;
        }

        .stat {
            text-align: center;
            padding: 16px;
            background-color: #1c1e21;
            border: 1px solid #2f3336;
            border-radius: 8px;
            transition: transform 0.2s, box-shadow 0.2s;
        }

        .stat:hover {
            transform: translateY(-2px);
            box-shadow: 0 2px 8px rgba(255, 255, 255, 0.1);
        }

        .stat-value {
            font-size: 32px;
            font-weight: 700;
            color: #e4e6eb;
            margin-bottom: 4px;
        }

        .stat-label {
            font-size: 13px;
            color: #8b949e;
            text-transform: uppercase;
            font-weight: 500;
            letter-spacing: 0.5px;
        }

        .feed {
            padding: 24px 0 80px;
        }

        .entry {
            padding: 0;
            margin-bottom: 24px;
            background-color: #1c1e21;
            border: 1px solid #2f3336;
            border-radius: 8px;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            display: flex;
            flex-direction: column;
        }

        .entry:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(255, 255, 255, 0.15);
        }

        .entry-header {
            display: flex;
            align-items: center;
            padding: 14px 16px;
            border-bottom: 1px solid #2f3336;
        }

        .avatar {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            background: #ffffff;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 600;
            font-size: 14px;
            color: #ffffff;
            margin-right: 12px;
            flex-shrink: 0;
            padding: 2px;
        }

        .avatar-inner {
            width: 100%;
            height: 100%;
            background: #1c1e21;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #ffffff;
        }

        .entry-info {
            flex: 1;
        }

        .username {
            display: none;
        }

        .timestamp {
            color: #8b949e;
            font-size: 12px;
            font-weight: 400;
            margin-top: 4px;
            display: block;
        }

        .entry-photo-container {
            width: 100%;
            max-width: 1024px;
            height: 0;
            padding-bottom: 100%;
            position: relative;
            overflow: hidden;
            background-color: #000000;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .entry-photo {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            max-width: 100%;
            max-height: 100%;
            width: auto;
            height: auto;
            object-fit: contain;
            display: block;
        }

        .entry-actions {
            display: none;
        }

        .action-btn {
            background: none;
            border: none;
            cursor: pointer;
            font-size: 24px;
            padding: 8px;
            line-height: 1;
            color: #e4e6eb;
            transition: opacity 0.2s;
        }

        .action-btn:hover {
            opacity: 0.6;
        }

        .entry-content {
            font-size: 15px;
            color: #e4e6eb;
            line-height: 20px;
            padding: 16px 16px 4px;
            word-wrap: break-word;
            overflow-wrap: break-word;
        }

        .entry-content a {
            color: #1d9bf0;
            text-decoration: none;
            font-weight: 500;
        }

        .entry-content a:hover {
            text-decoration: underline;
        }

        .entry-content-username {
            display: none;
        }

        .entry-timestamp {
            padding: 4px 16px 16px;
            color: #8b949e;
            font-size: 10px;
            text-transform: uppercase;
            letter-spacing: 0.2px;
        }

        .entry-footer {
            display: none;
        }

        .entry-id {
            color: #8b949e;
            font-size: 12px;
            font-weight: 500;
        }

        .entry-badge {
            background-color: #1d9bf0;
            color: #ffffff;
            padding: 4px 12px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .empty-state {
            padding: 80px 32px;
            text-align: center;
        }

        .empty-state-icon {
            width: 120px;
            height: 120px;
            margin: 0 auto 24px;
            background: #1c1e21;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 64px;
            border: 2px solid #2f3336;
        }

        .empty-state-title {
            font-size: 28px;
            font-weight: 600;
            margin-bottom: 12px;
            color: #e4e6eb;
        }

        .empty-state-text {
            font-size: 16px;
            color: #8b949e;
            max-width: 400px;
            margin: 0 auto;
        }

        .loading-indicator {
            text-align: center;
            padding: 24px;
            color: #8b949e;
            font-size: 14px;
            display: none;
        }

        .loading-indicator.show {
            display: block;
        }

        .loading-spinner {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 2px solid #2f3336;
            border-top-color: #1d9bf0;
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
            margin-right: 8px;
            vertical-align: middle;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        .end-message {
            text-align: center;
            padding: 24px;
            color: #8b949e;
            font-size: 14px;
            display: none;
        }

        .end-message.show {
            display: block;
        }

        @media (max-width: 768px) {
            .container {
                background-color: #1c1e21;
            }

            .header {
                padding: 16px 20px;
            }

            .header h1 {
                font-size: 24px;
            }

            .stats {
                padding: 12px 16px;
                gap: 12px;
            }

            .stat {
                padding: 12px;
            }

            .stat-value {
                font-size: 24px;
            }

            .feed {
                padding: 0 0 80px;
            }

            .entry {
                margin-bottom: 12px;
                border-radius: 0;
                border-left: none;
                border-right: none;
            }
        }
`

const viewerTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.SiteTitle}}</title>
    <style>
{{.ThemeCSS}}
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="header-content">
                <div>
                    <h1>{{.SiteTitle}}</h1>
                    <div class="subtitle">{{.SiteSubtitle}}</div>
                </div>
            </div>
        </div>

        <div class="feed" id="feed">
            {{if .Entries}}
                {{range .Entries}}
                <div class="entry" data-slug="{{.Slug}}">
                    <div class="entry-header">
                        <div class="avatar">
                            {{if and (eq $.AvatarPreference "avatar") $.AvatarPath}}
                                <img src="/uploads/{{$.AvatarPath}}" alt="Avatar" style="width: 100%; height: 100%; border-radius: 50%; object-fit: cover;">
                            {{else}}
                                <div class="avatar-inner">{{$.UserInitial}}</div>
                            {{end}}
                        </div>
                    </div>
                    {{if .HasPhoto}}
                    <div class="entry-photo-container">
                        <img src="{{.Photo}}" alt="Entry photo" class="entry-photo">
                    </div>
                    {{end}}
                    {{if .HasAudio}}
                    <div class="entry-photo-container">
                        {{if .HasThumbnail}}
                        <img src="{{.Thumbnail}}" alt="Audio cover" class="entry-photo" style="object-fit: cover;">
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
                                <source src="{{.Photo}}">
                                Your browser does not support the audio element.
                            </audio>
                        </div>
                    </div>
                    {{end}}
                    {{if .HasVideo}}
                    <div class="entry-photo-container">
                        {{if .HasThumbnail}}
                        <div class="video-thumbnail-wrapper" onclick="this.style.display='none'; this.nextElementSibling.style.display='block'; this.nextElementSibling.querySelector('video').play();">
                            <img src="{{.Thumbnail}}" alt="Video thumbnail" class="entry-photo" style="object-fit: cover; cursor: pointer;">
                            <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); background: rgba(0,0,0,0.6); border-radius: 50%; width: 80px; height: 80px; display: flex; align-items: center; justify-content: center; cursor: pointer;">
                                <svg width="40" height="40" viewBox="0 0 24 24" fill="white">
                                    <polygon points="5 3 19 12 5 21 5 3"></polygon>
                                </svg>
                            </div>
                        </div>
                        <div style="display: none;">
                            <video controls class="entry-photo" style="object-fit: contain;">
                                <source src="{{.Photo}}" type="video/mp4">
                                Your browser does not support the video element.
                            </video>
                        </div>
                        {{else}}
                        <video controls class="entry-photo" style="object-fit: contain;">
                            <source src="{{.Photo}}" type="video/mp4">
                            Your browser does not support the video element.
                        </video>
                        {{end}}
                    </div>
                    {{end}}
                    <div class="entry-content">{{.Content}}</div>
                    <div class="entry-timestamp">{{.TimeAgo}}</div>
                </div>
                {{end}}
            {{else}}
                <div class="empty-state">
                    <div class="empty-state-text">Your story starts here.</div>
                </div>
            {{end}}
        </div>

        <div class="loading-indicator" id="loading">
            <span class="loading-spinner"></span>
            Loading more entries...
        </div>

        <div class="end-message" id="endMessage">
            No more entries to load
        </div>
    </div>

    <script>
        let currentOffset = {{.InitialCount}};
        let isLoading = false;
        let hasMore = {{.HasMore}};
        const userInitial = '{{.UserInitial}}';
        const avatarPath = '{{.AvatarPath}}';
        const avatarPreference = '{{.AvatarPreference}}';

        // Add click handlers to all entries
        document.addEventListener('DOMContentLoaded', function() {
            addClickHandlersToEntries();
        });

        function addClickHandlersToEntries() {
            const entries = document.querySelectorAll('.entry[data-slug]');
            entries.forEach(function(entry) {
                if (!entry.hasAttribute('data-click-added')) {
                    entry.setAttribute('data-click-added', 'true');
                    entry.addEventListener('click', function(e) {
                        const slug = this.getAttribute('data-slug');
                        if (slug) {
                            window.location.href = '/posts/' + slug + '/';
                        }
                    });
                }
            });
        }

        function timeAgo(timestamp) {
            const now = new Date();
            const past = new Date(timestamp);
            const seconds = Math.floor((now - past) / 1000);

            if (seconds < 60) {
                return seconds <= 1 ? 'just now' : seconds + 's ago';
            }

            const minutes = Math.floor(seconds / 60);
            if (minutes < 60) {
                return minutes === 1 ? '1 minute ago' : minutes + ' minutes ago';
            }

            const hours = Math.floor(minutes / 60);
            if (hours < 24) {
                return hours === 1 ? '1 hour ago' : hours + ' hours ago';
            }

            const days = Math.floor(hours / 24);
            return days === 1 ? '1 day ago' : days + ' days ago';
        }

        function createEntryElement(entry) {
            const entryDiv = document.createElement('div');
            entryDiv.className = 'entry';

            // Add slug as data attribute for click handler
            if (entry.Slug) {
                entryDiv.setAttribute('data-slug', entry.Slug);
            }

            let mediaHtml = '';
            if (entry.HasPhoto && entry.Photo) {
                mediaHtml = '<div class="entry-photo-container"><img src="' + entry.Photo + '" alt="Entry photo" class="entry-photo"></div>';
            } else if (entry.HasAudio && entry.Photo) {
                if (entry.HasThumbnail && entry.Thumbnail) {
                    mediaHtml = '<div class="entry-photo-container"><img src="' + entry.Thumbnail + '" alt="Audio cover" class="entry-photo" style="object-fit: cover;"><div style="position: absolute; bottom: 0; left: 0; right: 0; padding: 16px; background: linear-gradient(transparent, rgba(0,0,0,0.8));"><audio controls style="width: 100%; filter: invert(1) hue-rotate(180deg);"><source src="' + entry.Photo + '">Your browser does not support the audio element.</audio></div></div>';
                } else {
                    mediaHtml = '<div class="entry-photo-container"><div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); text-align: center;"><svg width="80" height="80" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="opacity: 0.7;"><path d="M9 18V5l12-2v13"></path><circle cx="6" cy="18" r="3"></circle><circle cx="18" cy="16" r="3"></circle></svg></div><div style="position: absolute; bottom: 0; left: 0; right: 0; padding: 16px; background: linear-gradient(transparent, rgba(0,0,0,0.8));"><audio controls style="width: 100%; filter: invert(1) hue-rotate(180deg);"><source src="' + entry.Photo + '">Your browser does not support the audio element.</audio></div></div>';
                }
            } else if (entry.HasVideo && entry.Photo) {
                if (entry.HasThumbnail && entry.Thumbnail) {
                    const uniqueId = 'video-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9);
                    mediaHtml = '<div class="entry-photo-container"><div class="video-thumbnail-wrapper" onclick="this.style.display=\'none\'; this.nextElementSibling.style.display=\'block\'; this.nextElementSibling.querySelector(\'video\').play();"><img src="' + entry.Thumbnail + '" alt="Video thumbnail" class="entry-photo" style="object-fit: cover; cursor: pointer;"><div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); background: rgba(0,0,0,0.6); border-radius: 50%; width: 80px; height: 80px; display: flex; align-items: center; justify-content: center; cursor: pointer;"><svg width="40" height="40" viewBox="0 0 24 24" fill="white"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg></div></div><div style="display: none;"><video controls class="entry-photo" style="object-fit: contain;"><source src="' + entry.Photo + '" type="video/mp4">Your browser does not support the video element.</video></div></div>';
                } else {
                    mediaHtml = '<div class="entry-photo-container"><video controls class="entry-photo" style="object-fit: contain;"><source src="' + entry.Photo + '" type="video/mp4">Your browser does not support the video element.</video></div>';
                }
            }

            let avatarHtml = (avatarPreference === 'avatar' && avatarPath)
                ? '<img src="/uploads/' + avatarPath + '" alt="Avatar" style="width: 100%; height: 100%; border-radius: 50%; object-fit: cover;">'
                : '<div class="avatar-inner">' + userInitial + '</div>';

            entryDiv.innerHTML = '<div class="entry-header"><div class="avatar">' + avatarHtml + '</div></div>' +
                mediaHtml +
                '<div class="entry-content">' + entry.Content + '</div>' +
                '<div class="entry-timestamp">' + entry.TimeAgo + '</div>';

            return entryDiv;
        }

        async function loadMoreEntries() {
            if (isLoading || !hasMore) return;

            isLoading = true;
            const loadingEl = document.getElementById('loading');
            loadingEl.classList.add('show');

            try {
                const response = await fetch('/api/entries?offset=' + currentOffset + '&limit=10');
                const data = await response.json();

                if (data.entries && data.entries.length > 0) {
                    const feed = document.getElementById('feed');

                    data.entries.forEach(entry => {
                        feed.appendChild(createEntryElement(entry));
                    });

                    // Add click handlers to newly loaded entries
                    addClickHandlersToEntries();

                    currentOffset += data.entries.length;
                    hasMore = data.hasMore;

                    if (!hasMore) {
                        document.getElementById('endMessage').classList.add('show');
                    }
                } else {
                    hasMore = false;
                    document.getElementById('endMessage').classList.add('show');
                }
            } catch (error) {
                console.error('Error loading entries:', error);
            } finally {
                isLoading = false;
                loadingEl.classList.remove('show');
            }
        }

        // Infinite scroll
        window.addEventListener('scroll', () => {
            if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight - 500) {
                loadMoreEntries();
            }
        });
    </script>
</body>
</html>`
