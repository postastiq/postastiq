package main

const notFoundTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>404 - Page Not Found | {{.SiteTitle}}</title>
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
            <div style="text-align: center; padding: 80px 20px;">
                <div style="font-size: 120px; font-weight: 700; color: #dbdbdb; line-height: 1;">404</div>
                <div style="font-size: 24px; font-weight: 600; color: #262626; margin-top: 20px;">Page Not Found</div>
                <div style="font-size: 14px; color: #8e8e8e; margin-top: 12px; max-width: 400px; margin-left: auto; margin-right: auto;">
                    Sorry, the page you're looking for doesn't exist or has been moved.
                </div>
                <div style="margin-top: 32px;">
                    <a href="/" style="display: inline-block; background-color: #0095f6; color: #ffffff; padding: 12px 32px; border-radius: 8px; text-decoration: none; font-size: 14px; font-weight: 600; transition: background-color 0.2s;">
                        Go to Home
                    </a>
                </div>
            </div>

            {{if .RecentEntries}}
            <div style="margin-top: 40px; padding-top: 40px; border-top: 1px solid #dbdbdb;">
                <div style="text-align: center; font-size: 18px; font-weight: 600; color: #262626; margin-bottom: 24px;">
                    Recent Posts
                </div>
                {{range .RecentEntries}}
                <div class="entry" data-slug="{{.Slug}}" style="margin-bottom: 16px;">
                    <div class="entry-header">
                        <div class="avatar">
                            <div class="avatar-inner">{{$.UserInitial}}</div>
                        </div>
                    </div>
                    {{if .HasPhoto}}
                    <div class="entry-photo-container">
                        <img src="{{.Photo}}" alt="Entry photo" class="entry-photo">
                    </div>
                    {{end}}
                    <div class="entry-content">{{.Content}}</div>
                    <div class="entry-timestamp">{{.TimeAgo}}</div>
                </div>
                {{end}}
            </div>
            {{end}}
        </div>
    </div>

    <script>
        // Add click handlers to entries
        document.addEventListener('DOMContentLoaded', function() {
            const entries = document.querySelectorAll('.entry[data-slug]');
            entries.forEach(function(entry) {
                entry.addEventListener('click', function(e) {
                    const slug = this.getAttribute('data-slug');
                    if (slug) {
                        window.location.href = '/posts/' + slug + '/';
                    }
                });
            });
        });
    </script>
</body>
</html>`
