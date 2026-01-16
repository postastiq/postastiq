package main

const defaultThemeCSS = `
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: #ffffff;
            color: #1a1a1a;
            line-height: 1.6;
        }

        .container {
            max-width: 1100px;
            margin: 0 auto;
            padding: 0 32px;
            background-color: #ffffff;
            min-height: 100vh;
        }

        .header {
            background: #ffffff;
            padding: 32px 0;
            color: #1a1a1a;
            border-bottom: none;
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
            font-size: 24px;
            font-weight: 700;
            display: flex;
            align-items: center;
            gap: 12px;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            letter-spacing: -0.02em;
        }

        .subtitle {
            font-size: 14px;
            color: #6b7280;
            margin-top: 2px;
            font-weight: 400;
        }

        .stats {
            padding: 16px 0;
            background-color: #ffffff;
            display: none;
        }

        .stat {
            text-align: center;
            padding: 16px;
            background-color: #ffffff;
            border: 1px solid #e5e7eb;
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
            color: #1a1a1a;
            margin-bottom: 4px;
        }

        .stat-label {
            font-size: 13px;
            color: #6b7280;
            text-transform: uppercase;
            font-weight: 500;
            letter-spacing: 0.5px;
        }

        .feed {
            padding: 0 0 80px;
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 32px;
        }

        .entry {
            padding: 0;
            background-color: transparent;
            border: none;
            border-radius: 0;
            cursor: pointer;
            display: flex;
            flex-direction: column;
            gap: 12px;
        }

        .entry:hover .entry-photo-container {
            transform: scale(1.015);
        }

        .entry:hover .entry-content {
            opacity: 0.8;
        }

        /* Featured first post - spans 2 columns */
        .entry:first-child {
            grid-column: span 2;
        }

        .entry:first-child .entry-photo-container {
            aspect-ratio: 16 / 9;
            padding-bottom: 0;
            height: auto;
        }

        .entry:first-child .entry-content {
            font-size: 24px;
            font-weight: 600;
            line-height: 1.2;
            letter-spacing: -0.02em;
        }

        .entry:first-child .entry-timestamp {
            font-size: 14px;
        }

        .entry-header {
            display: none;
        }

        .avatar {
            display: none;
        }

        .avatar-inner {
            display: none;
        }

        .entry-info {
            flex: 1;
        }

        .username {
            display: none;
        }

        .timestamp {
            color: #6b7280;
            font-size: 13px;
            font-weight: 400;
            display: block;
            order: -1;
        }

        .entry-photo-container {
            width: 100%;
            max-width: none;
            height: auto;
            padding-bottom: 75%;
            aspect-ratio: 4 / 3;
            position: relative;
            overflow: hidden;
            background-color: #f3f4f6;
            border-radius: 18px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }

        .entry-photo-container:hover {
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
        }

        .entry-photo {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            object-fit: cover;
            display: block;
            transform: none;
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
            color: #1a1a1a;
            transition: opacity 0.2s;
        }

        .action-btn:hover {
            opacity: 0.6;
        }

        .entry-content {
            font-size: 18px;
            font-weight: 500;
            color: #1a1a1a;
            line-height: 1.25;
            padding: 0;
            word-wrap: break-word;
            overflow-wrap: break-word;
            transition: opacity 0.2s ease;
            display: -webkit-box;
            -webkit-line-clamp: 3;
            -webkit-box-orient: vertical;
            overflow: hidden;
            order: 1;
        }

        .entry-content a {
            color: #1a1a1a;
            text-decoration: none;
            font-weight: 600;
        }

        .entry-content a:hover {
            text-decoration: none;
        }

        .entry-content-username {
            display: none;
        }

        .entry-timestamp {
            padding: 0;
            color: #6b7280;
            font-size: 13px;
            text-transform: none;
            letter-spacing: 0.01em;
            order: 0;
            font-weight: 400;
        }

        .entry-footer {
            display: none;
        }

        .entry-id {
            color: #6b7280;
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
            grid-column: 1 / -1;
        }

        .empty-state-icon {
            width: 120px;
            height: 120px;
            margin: 0 auto 24px;
            background: #f3f4f6;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 64px;
            border: 2px solid #e5e7eb;
        }

        .empty-state-title {
            font-size: 28px;
            font-weight: 600;
            margin-bottom: 12px;
            color: #1a1a1a;
        }

        .empty-state-text {
            font-size: 16px;
            color: #6b7280;
            max-width: 400px;
            margin: 0 auto;
        }

        .loading-indicator {
            text-align: center;
            padding: 24px;
            color: #6b7280;
            font-size: 14px;
            display: none;
            grid-column: 1 / -1;
        }

        .loading-indicator.show {
            display: block;
        }

        .loading-spinner {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 2px solid #e5e7eb;
            border-top-color: #1a1a1a;
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
            color: #6b7280;
            font-size: 14px;
            display: none;
            grid-column: 1 / -1;
        }

        .end-message.show {
            display: block;
        }

        /* Tablet: 2 columns */
        @media (max-width: 900px) {
            .container {
                padding: 0 24px;
            }

            .feed {
                grid-template-columns: repeat(2, 1fr);
                gap: 24px;
            }

            .entry:first-child {
                grid-column: span 2;
            }

            .entry:first-child .entry-content {
                font-size: 22px;
            }

            .entry-content {
                font-size: 17px;
            }
        }

        /* Mobile: 1 column */
        @media (max-width: 600px) {
            .container {
                padding: 0 16px;
                background-color: #ffffff;
            }

            .header {
                padding: 24px 0;
            }

            .header h1 {
                font-size: 20px;
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
                grid-template-columns: 1fr;
                gap: 28px;
                padding: 0 0 80px;
            }

            .entry:first-child {
                grid-column: span 1;
            }

            .entry:first-child .entry-content {
                font-size: 18px;
            }

            .entry:first-child .entry-timestamp {
                font-size: 13px;
            }

            .entry-photo-container {
                border-radius: 14px;
            }

            .entry-content {
                font-size: 16px;
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
            color: #f0f0f0;
            line-height: 1.6;
        }

        .container {
            max-width: 1100px;
            margin: 0 auto;
            padding: 0 32px;
            background-color: #0a0a0a;
            min-height: 100vh;
        }

        .header {
            background: #0a0a0a;
            padding: 32px 0;
            color: #f0f0f0;
            border-bottom: none;
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
            font-size: 24px;
            font-weight: 700;
            display: flex;
            align-items: center;
            gap: 12px;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            letter-spacing: -0.02em;
        }

        .subtitle {
            font-size: 14px;
            color: #9ca3af;
            margin-top: 2px;
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
            background-color: #1a1a1a;
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
            color: #f0f0f0;
            margin-bottom: 4px;
        }

        .stat-label {
            font-size: 13px;
            color: #9ca3af;
            text-transform: uppercase;
            font-weight: 500;
            letter-spacing: 0.5px;
        }

        .feed {
            padding: 0 0 80px;
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 32px;
        }

        .entry {
            padding: 0;
            background-color: transparent;
            border: none;
            border-radius: 0;
            cursor: pointer;
            display: flex;
            flex-direction: column;
            gap: 12px;
        }

        .entry:hover .entry-photo-container {
            transform: scale(1.015);
        }

        .entry:hover .entry-content {
            opacity: 0.8;
        }

        /* Featured first post - spans 2 columns */
        .entry:first-child {
            grid-column: span 2;
        }

        .entry:first-child .entry-photo-container {
            aspect-ratio: 16 / 9;
            padding-bottom: 0;
            height: auto;
        }

        .entry:first-child .entry-content {
            font-size: 24px;
            font-weight: 600;
            line-height: 1.2;
            letter-spacing: -0.02em;
        }

        .entry:first-child .entry-timestamp {
            font-size: 14px;
        }

        .entry-header {
            display: none;
        }

        .avatar {
            display: none;
        }

        .avatar-inner {
            display: none;
        }

        .entry-info {
            flex: 1;
        }

        .username {
            display: none;
        }

        .timestamp {
            color: #9ca3af;
            font-size: 13px;
            font-weight: 400;
            display: block;
            order: -1;
        }

        .entry-photo-container {
            width: 100%;
            max-width: none;
            height: auto;
            padding-bottom: 75%;
            aspect-ratio: 4 / 3;
            position: relative;
            overflow: hidden;
            background-color: #1a1a1a;
            border-radius: 18px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }

        .entry-photo-container:hover {
            box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
        }

        .entry-photo {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            object-fit: cover;
            display: block;
            transform: none;
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
            color: #f0f0f0;
            transition: opacity 0.2s;
        }

        .action-btn:hover {
            opacity: 0.6;
        }

        .entry-content {
            font-size: 18px;
            font-weight: 500;
            color: #f0f0f0;
            line-height: 1.25;
            padding: 0;
            word-wrap: break-word;
            overflow-wrap: break-word;
            transition: opacity 0.2s ease;
            display: -webkit-box;
            -webkit-line-clamp: 3;
            -webkit-box-orient: vertical;
            overflow: hidden;
            order: 1;
        }

        .entry-content a {
            color: #f0f0f0;
            text-decoration: none;
            font-weight: 600;
        }

        .entry-content a:hover {
            text-decoration: none;
        }

        .entry-content-username {
            display: none;
        }

        .entry-timestamp {
            padding: 0;
            color: #9ca3af;
            font-size: 13px;
            text-transform: none;
            letter-spacing: 0.01em;
            order: 0;
            font-weight: 400;
        }

        .entry-footer {
            display: none;
        }

        .entry-id {
            color: #9ca3af;
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
            grid-column: 1 / -1;
        }

        .empty-state-icon {
            width: 120px;
            height: 120px;
            margin: 0 auto 24px;
            background: #1a1a1a;
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
            color: #f0f0f0;
        }

        .empty-state-text {
            font-size: 16px;
            color: #9ca3af;
            max-width: 400px;
            margin: 0 auto;
        }

        .loading-indicator {
            text-align: center;
            padding: 24px;
            color: #9ca3af;
            font-size: 14px;
            display: none;
            grid-column: 1 / -1;
        }

        .loading-indicator.show {
            display: block;
        }

        .loading-spinner {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 2px solid #2f3336;
            border-top-color: #f0f0f0;
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
            color: #9ca3af;
            font-size: 14px;
            display: none;
            grid-column: 1 / -1;
        }

        .end-message.show {
            display: block;
        }

        /* Tablet: 2 columns */
        @media (max-width: 900px) {
            .container {
                padding: 0 24px;
            }

            .feed {
                grid-template-columns: repeat(2, 1fr);
                gap: 24px;
            }

            .entry:first-child {
                grid-column: span 2;
            }

            .entry:first-child .entry-content {
                font-size: 22px;
            }

            .entry-content {
                font-size: 17px;
            }
        }

        /* Mobile: 1 column */
        @media (max-width: 600px) {
            .container {
                padding: 0 16px;
                background-color: #0a0a0a;
            }

            .header {
                padding: 24px 0;
            }

            .header h1 {
                font-size: 20px;
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
                grid-template-columns: 1fr;
                gap: 28px;
                padding: 0 0 80px;
            }

            .entry:first-child {
                grid-column: span 1;
            }

            .entry:first-child .entry-content {
                font-size: 18px;
            }

            .entry:first-child .entry-timestamp {
                font-size: 13px;
            }

            .entry-photo-container {
                border-radius: 14px;
            }

            .entry-content {
                font-size: 16px;
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
                    {{if .EnableSubtitle}}<div class="subtitle">{{.SiteSubtitle}}</div>{{end}}
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
