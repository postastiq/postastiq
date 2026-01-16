package main

const adminTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin - {{.PageTitle}}</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: #fafafa;
            color: #262626;
            line-height: 1.6;
            min-height: 100vh;
        }

        /* Layout */
        .admin-layout {
            display: flex;
            min-height: 100vh;
        }

        /* Sidebar */
        .sidebar {
            width: 240px;
            background-color: #ffffff;
            border-right: 1px solid #dbdbdb;
            padding: 20px 0;
            position: fixed;
            height: 100vh;
            overflow-y: auto;
            z-index: 100;
        }
        .sidebar-header {
            padding: 0 20px 20px 20px;
            border-bottom: 1px solid #efefef;
            margin-bottom: 20px;
        }
        .sidebar-logo {
            font-size: 20px;
            font-weight: 700;
            color: #262626;
            text-decoration: none;
        }
        .sidebar-section {
            margin-bottom: 24px;
        }
        .sidebar-section-title {
            font-size: 11px;
            font-weight: 600;
            color: #8e8e8e;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            padding: 0 20px;
            margin-bottom: 8px;
        }
        .sidebar-nav {
            list-style: none;
        }
        .sidebar-nav a {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 10px 20px;
            color: #262626;
            text-decoration: none;
            font-size: 14px;
            font-weight: 500;
            transition: background-color 0.2s;
        }
        .sidebar-nav a:hover {
            background-color: #f0f0f0;
        }
        .sidebar-nav a.active {
            background-color: #efefef;
            font-weight: 600;
        }
        .sidebar-nav a svg {
            width: 18px;
            height: 18px;
            stroke: currentColor;
            fill: none;
            stroke-width: 2;
            stroke-linecap: round;
            stroke-linejoin: round;
        }
        .sidebar-footer {
            position: absolute;
            bottom: 0;
            left: 0;
            right: 0;
            padding: 20px;
            border-top: 1px solid #efefef;
            background: #ffffff;
        }
        .sidebar-footer a {
            display: flex;
            align-items: center;
            gap: 12px;
            color: #8e8e8e;
            text-decoration: none;
            font-size: 14px;
            font-weight: 500;
            padding: 8px 0;
        }
        .sidebar-footer a:hover {
            color: #262626;
        }

        /* Main Content */
        .main-content {
            flex: 1;
            margin-left: 240px;
            padding: 24px;
            min-height: 100vh;
        }
        .content-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 24px;
            padding-bottom: 16px;
            border-bottom: 1px solid #dbdbdb;
        }
        .content-header h1 {
            font-size: 24px;
            font-weight: 600;
            color: #262626;
        }
        .content-container {
            max-width: 700px;
            background-color: #ffffff;
            border: 1px solid #dbdbdb;
            border-radius: 8px;
            padding: 24px;
        }

        /* Messages */
        .message { padding: 12px 16px; margin-bottom: 20px; border-radius: 8px; font-size: 14px; line-height: 18px; }
        .success { background-color: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .error { background-color: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }

        /* Forms */
        .form-group { margin-bottom: 20px; }
        label { display: block; margin-bottom: 8px; font-weight: 600; font-size: 14px; color: #262626; }
        input[type="file"] {
            width: 100%;
            padding: 10px 12px;
            border: 1px solid #dbdbdb;
            border-radius: 8px;
            font-size: 14px;
            background-color: #fafafa;
            cursor: pointer;
        }
        input[type="text"],
        input[type="password"] {
            width: 100%;
            padding: 12px 16px;
            border: 1px solid #dbdbdb;
            border-radius: 8px;
            font-size: 14px;
            font-family: inherit;
            background-color: #fafafa;
            transition: border-color 0.2s, background-color 0.2s;
        }
        input[type="text"]:hover,
        input[type="password"]:hover {
            border-color: #a8a8a8;
        }
        input[type="text"]:focus,
        input[type="password"]:focus {
            outline: none;
            border-color: #0095f6;
            background-color: #ffffff;
        }
        textarea {
            width: 100%;
            padding: 12px 16px;
            border: 1px solid #dbdbdb;
            border-radius: 8px;
            font-size: 14px;
            font-family: inherit;
            resize: vertical;
            min-height: 120px;
            background-color: #fafafa;
            transition: border-color 0.2s, background-color 0.2s;
        }
        textarea:hover {
            border-color: #a8a8a8;
        }
        textarea:focus { outline: none; border-color: #0095f6; background-color: #ffffff; }
        .file-info { font-size: 12px; color: #8e8e8e; margin-top: 6px; }
        .custom-file-upload {
            display: inline-flex;
            align-items: center;
            padding: 12px 20px;
            background-color: #f0f0f0;
            color: #262626;
            border: 2px dashed #dbdbdb;
            border-radius: 8px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.3s ease;
            width: 100%;
            justify-content: center;
        }
        .custom-file-upload:hover {
            background-color: #e8e8e8;
            border-color: #0095f6;
        }
        .custom-file-upload svg {
            flex-shrink: 0;
        }
        .preview-container { margin-top: 12px; display: none; padding: 12px; background-color: #fafafa; border-radius: 8px; }
        .preview-container.active { display: block; }
        .preview-image { max-width: 100%; height: auto; border-radius: 4px; display: block; }

        /* Buttons */
        button, .btn {
            background-color: #000000;
            color: #ffffff;
            padding: 12px 24px;
            border: none;
            cursor: pointer;
            border-radius: 8px;
            font-size: 14px;
            font-weight: 600;
            transition: transform 0.2s;
            text-decoration: none;
            display: inline-block;
        }
        button:hover, .btn:hover { transform: translateY(-1px); }
        button.full-width, .btn.full-width { width: 100%; }
        button.btn-danger, .btn.btn-danger { background-color: #ed4956; }
        button.btn-secondary, .btn.btn-secondary { background-color: #8e8e8e; }
        button.btn-primary, .btn.btn-primary { background-color: #0095f6; }

        /* Entry Cards */
        .entry-card {
            background-color: #ffffff;
            border: 1px solid #dbdbdb;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 16px;
        }
        .entry-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 12px;
        }
        .entry-meta {
            font-size: 13px;
            color: #8e8e8e;
        }
        .entry-title {
            font-weight: 600;
            font-size: 16px;
            margin-bottom: 8px;
            color: #262626;
        }
        .entry-content {
            margin-bottom: 12px;
            font-size: 14px;
            white-space: pre-wrap;
            color: #555;
            max-height: 80px;
            overflow: hidden;
        }
        .entry-media-badge {
            display: inline-flex;
            align-items: center;
            gap: 4px;
            font-size: 12px;
            color: #8e8e8e;
            background: #f0f0f0;
            padding: 4px 8px;
            border-radius: 4px;
            margin-bottom: 12px;
        }
        .entry-photo {
            max-width: 150px;
            max-height: 100px;
            border-radius: 4px;
            margin-bottom: 12px;
            cursor: pointer;
            object-fit: cover;
        }
        .entry-actions { display: flex; gap: 8px; }
        .entry-actions button { flex: 1; padding: 8px 16px; font-size: 13px; }

        /* Pagination */
        .pagination {
            display: flex;
            justify-content: center;
            align-items: center;
            gap: 8px;
            margin-top: 24px;
            padding-top: 24px;
            border-top: 1px solid #efefef;
        }
        .pagination a, .pagination span {
            padding: 8px 14px;
            border: 1px solid #dbdbdb;
            border-radius: 6px;
            text-decoration: none;
            color: #262626;
            font-size: 14px;
            font-weight: 500;
        }
        .pagination a:hover {
            background-color: #f0f0f0;
        }
        .pagination .active {
            background-color: #000000;
            color: #ffffff;
            border-color: #000000;
        }
        .pagination .disabled {
            color: #8e8e8e;
            cursor: not-allowed;
        }
        .pagination-info {
            font-size: 13px;
            color: #8e8e8e;
            margin-bottom: 16px;
        }

        /* Modal */
        .modal {
            display: none;
            position: fixed;
            z-index: 1000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0, 0, 0, 0.5);
            overflow-y: auto;
        }
        .modal.active { display: flex; align-items: center; justify-content: center; padding: 20px; }
        .modal-content {
            background-color: #ffffff;
            border-radius: 8px;
            max-width: 500px;
            width: 100%;
            padding: 32px;
        }
        .modal-header { margin-bottom: 24px; }
        .modal-header h2 { font-size: 20px; font-weight: 600; }
        .modal-actions { display: flex; gap: 8px; margin-top: 20px; }
        .modal-actions button { flex: 1; }

        /* Audio Recording Styles */
        .audio-source-options { display: none; margin-top: 12px; padding: 12px; background-color: #fafafa; border-radius: 8px; }
        .audio-source-options.active { display: block; }
        .audio-source-radio { display: flex; gap: 20px; margin-bottom: 12px; }
        .recording-ui { display: none; margin-top: 12px; }
        .recording-ui.active { display: block; }
        .record-btn {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            padding: 12px 24px;
            border-radius: 8px;
            font-size: 14px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.2s;
            border: none;
        }
        .record-btn.start { background-color: #ed4956; color: white; }
        .record-btn.start:hover { background-color: #dc3545; }
        .record-btn.stop { background-color: #8e8e8e; color: white; }
        .record-btn.stop:hover { background-color: #6c757d; }
        .recording-indicator {
            display: none;
            align-items: center;
            gap: 8px;
            margin-top: 12px;
            padding: 12px;
            background-color: #fff3cd;
            border-radius: 8px;
            font-size: 14px;
        }
        .recording-indicator.active { display: flex; }
        .recording-dot {
            width: 12px;
            height: 12px;
            background-color: #ed4956;
            border-radius: 50%;
            animation: pulse 1s infinite;
        }
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }
        .audio-preview { display: none; margin-top: 12px; padding: 12px; background-color: #d4edda; border-radius: 8px; }
        .audio-preview.active { display: block; }
        .audio-preview audio { width: 100%; margin-bottom: 8px; }
        .audio-preview-actions { display: flex; gap: 8px; }
        .audio-preview-actions button { flex: 1; padding: 8px 16px; font-size: 13px; }

        .empty-state { text-align: center; padding: 60px 20px; color: #8e8e8e; }
        .empty-state svg { width: 64px; height: 64px; stroke: #dbdbdb; margin-bottom: 16px; }

        /* Mobile Responsive */
        .mobile-header {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            background: #ffffff;
            border-bottom: 1px solid #dbdbdb;
            padding: 12px 16px;
            z-index: 101;
        }
        .mobile-header-content {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .hamburger-btn {
            background: none;
            border: none;
            padding: 8px;
            cursor: pointer;
        }
        .hamburger-btn svg {
            width: 24px;
            height: 24px;
            stroke: #262626;
        }
        .sidebar-overlay {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0,0,0,0.5);
            z-index: 99;
        }

        @media (max-width: 768px) {
            .mobile-header {
                display: block;
            }
            .sidebar {
                transform: translateX(-100%);
                transition: transform 0.3s ease;
            }
            .sidebar.open {
                transform: translateX(0);
            }
            .sidebar-overlay.open {
                display: block;
            }
            .main-content {
                margin-left: 0;
                padding: 80px 16px 16px 16px;
            }
            .content-container {
                padding: 16px;
            }
        }
    </style>
</head>
<body>
    <!-- Mobile Header -->
    <div class="mobile-header">
        <div class="mobile-header-content">
            <span class="sidebar-logo">Postastiq</span>
            <button class="hamburger-btn" onclick="toggleSidebar()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="3" y1="12" x2="21" y2="12"></line>
                    <line x1="3" y1="6" x2="21" y2="6"></line>
                    <line x1="3" y1="18" x2="21" y2="18"></line>
                </svg>
            </button>
        </div>
    </div>
    <div class="sidebar-overlay" onclick="toggleSidebar()"></div>

    <div class="admin-layout">
        <!-- Sidebar -->
        <aside class="sidebar" id="sidebar">
            <div class="sidebar-header">
                <a href="/admin" class="sidebar-logo">Postastiq</a>
            </div>

            <div class="sidebar-section">
                <div class="sidebar-section-title">Content</div>
                <ul class="sidebar-nav">
                    <li>
                        <a href="/admin" class="{{if eq .View "new"}}active{{end}}">
                            <svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="16"></line><line x1="8" y1="12" x2="16" y2="12"></line></svg>
                            New Post
                        </a>
                    </li>
                    <li>
                        <a href="/admin?view=posts" class="{{if eq .View "posts"}}active{{end}}">
                            <svg viewBox="0 0 24 24"><path d="M12 20h9M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
                            All Posts
                        </a>
                    </li>
                </ul>
            </div>

            <div class="sidebar-section">
                <div class="sidebar-section-title">Settings</div>
                <ul class="sidebar-nav">
                    <li>
                        <a href="/admin/settings">
                            <svg viewBox="0 0 24 24"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
                            Site Info
                        </a>
                    </li>
                    <li>
                        <a href="/admin/settings/appearance">
                            <svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line></svg>
                            Appearance
                        </a>
                    </li>
                    <li>
                        <a href="/admin/settings/security">
                            <svg viewBox="0 0 24 24"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
                            Security
                        </a>
                    </li>
                    <li>
                        <a href="/admin/settings/backup">
                            <svg viewBox="0 0 24 24"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>
                            Backup
                        </a>
                    </li>
                </ul>
            </div>

            <div class="sidebar-footer">
                <a href="/" target="_blank">
                    <svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" fill="none" stroke-width="2"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
                    View Site
                </a>
                <a href="/rss" target="_blank">
                    <svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" fill="none" stroke-width="2"><path d="M4 11a9 9 0 0 1 9 9"></path><path d="M4 4a16 16 0 0 1 16 16"></path><circle cx="5" cy="19" r="1"></circle></svg>
                    RSS Feed
                </a>
                <a href="/logout">
                    <svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" fill="none" stroke-width="2"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
                    Logout
                </a>
            </div>
        </aside>

        <!-- Main Content -->
        <main class="main-content">
            {{if .Message}}
            <div class="message {{.MessageType}}">{{.Message}}</div>
            {{end}}

            {{if eq .View "posts"}}
            <!-- Posts List View -->
            <div class="content-header">
                <h1>Posts</h1>
                <a href="/admin?view=new" class="btn btn-primary">+ New Post</a>
            </div>

            <div class="content-container">
            <div class="pagination-info">
                {{if .TotalEntries}}
                Showing {{.StartEntry}}-{{.EndEntry}} of {{.TotalEntries}} posts
                {{else}}
                No posts yet
                {{end}}
            </div>

            <div id="entriesContainer">
                {{if .Entries}}
                {{range .Entries}}
                <div class="entry-card">
                    <div class="entry-header">
                        <div>
                            <div class="entry-title">{{if .Title}}{{.Title}}{{else}}Untitled{{end}}</div>
                            <div class="entry-meta">{{.CreatedAt.Format "Jan 2, 2006 at 3:04 PM"}}</div>
                        </div>
                    </div>
                    <div class="entry-content">{{.Content}}</div>
                    {{if .PhotoPath}}
                    <div class="entry-media-badge">
                        {{if eq .MediaType "video"}}
                        <svg viewBox="0 0 24 24" width="14" height="14" stroke="currentColor" fill="none" stroke-width="2"><polygon points="23 7 16 12 23 17 23 7"></polygon><rect x="1" y="5" width="15" height="14" rx="2" ry="2"></rect></svg>
                        Video
                        {{else if eq .MediaType "audio"}}
                        <svg viewBox="0 0 24 24" width="14" height="14" stroke="currentColor" fill="none" stroke-width="2"><path d="M9 18V5l12-2v13"></path><circle cx="6" cy="18" r="3"></circle><circle cx="18" cy="16" r="3"></circle></svg>
                        Audio
                        {{else}}
                        <svg viewBox="0 0 24 24" width="14" height="14" stroke="currentColor" fill="none" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><circle cx="8.5" cy="8.5" r="1.5"></circle><polyline points="21 15 16 10 5 21"></polyline></svg>
                        Photo
                        {{end}}
                    </div>
                    {{if and .PhotoPath (or (eq .MediaType "photo") (eq .MediaType ""))}}
                    <img src="/uploads/{{.PhotoPath}}" alt="" class="entry-photo">
                    {{end}}
                    {{end}}
                    <div class="entry-actions">
                        <button onclick="openEditModal({{.ID}}, '{{jsEscape .Content}}')">Edit</button>
                        <button class="btn-danger" onclick="openDeleteModal({{.ID}})">Delete</button>
                    </div>
                </div>
                {{end}}
                {{else}}
                <div class="empty-state">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                        <path d="M12 20h9M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path>
                    </svg>
                    <p>No posts yet. Create your first post!</p>
                    <a href="/admin?view=new" class="btn" style="margin-top: 16px;">Create Post</a>
                </div>
                {{end}}
            </div>

            {{if and .Entries (gt .TotalPages 1)}}
            <div class="pagination">
                {{if gt .CurrentPage 1}}
                <a href="/admin?view=posts&page={{subtract .CurrentPage 1}}">Previous</a>
                {{else}}
                <span class="disabled">Previous</span>
                {{end}}

                {{range .PageNumbers}}
                {{if eq . $.CurrentPage}}
                <span class="active">{{.}}</span>
                {{else}}
                <a href="/admin?view=posts&page={{.}}">{{.}}</a>
                {{end}}
                {{end}}

                {{if lt .CurrentPage .TotalPages}}
                <a href="/admin?view=posts&page={{add .CurrentPage 1}}">Next</a>
                {{else}}
                <span class="disabled">Next</span>
                {{end}}
            </div>
            {{end}}
            </div>

            {{else if eq .View "new"}}
            <!-- New Post View -->
            <div class="content-header">
                <h1>New Post</h1>
            </div>

            <div class="content-container">
                <form method="POST" action="/admin/create" enctype="multipart/form-data" id="createForm">
                    <div class="form-group">
                        <label for="content">Content (max 2000 characters)</label>
                        <textarea name="content" id="content" maxlength="2000" placeholder="What's on your mind?" required></textarea>
                        <div class="file-info"><span id="charCount">0</span> / 2000 characters</div>
                    </div>

                    <div class="form-group">
                        <label>Media (optional)</label>
                        <div style="display: flex; gap: 20px; margin-top: 8px;">
                            <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                                <input type="radio" name="media_type" value="photo" checked class="mediaTypeRadio">
                                <span>Photo</span>
                            </label>
                            {{if .EnableAudioUploads}}
                            <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                                <input type="radio" name="media_type" value="audio" class="mediaTypeRadio">
                                <span>Audio</span>
                            </label>
                            {{end}}
                            <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                                <input type="radio" name="media_type" value="video" class="mediaTypeRadio">
                                <span>Video</span>
                            </label>
                        </div>
                    </div>

                    <div class="form-group" id="fileUploadGroup">
                        <input type="file" name="media" id="media" accept="image/*" style="display: none;">
                        <div class="custom-file-upload" id="fileUploadBtn" onclick="document.getElementById('media').click()">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 8px;">
                                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                                <polyline points="17 8 12 3 7 8"></polyline>
                                <line x1="12" y1="3" x2="12" y2="15"></line>
                            </svg>
                            <span id="mediaFileName">Choose file</span>
                        </div>
                        <div class="preview-container" id="photoPreview">
                            <img src="" alt="Preview" class="preview-image" id="previewImg">
                        </div>
                    </div>

                    <!-- Audio Source Options -->
                    {{if .EnableAudioUploads}}
                    <div class="audio-source-options" id="audioSourceOptions">
                        <div class="audio-source-radio">
                            <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                                <input type="radio" name="audio_source" value="upload" checked class="audioSourceRadio">
                                <span>Upload Audio File</span>
                            </label>
                            <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                                <input type="radio" name="audio_source" value="record" class="audioSourceRadio">
                                <span>Record Audio</span>
                            </label>
                        </div>

                        <div class="recording-ui" id="recordingUI">
                            <button type="button" class="record-btn start" id="startRecordBtn" onclick="startRecording()">
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                    <circle cx="12" cy="12" r="10"/>
                                </svg>
                                Start Recording
                            </button>
                            <button type="button" class="record-btn stop" id="stopRecordBtn" onclick="stopRecording()" style="display: none;">
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                    <rect x="6" y="6" width="12" height="12"/>
                                </svg>
                                Stop Recording
                            </button>

                            <div class="recording-indicator" id="recordingIndicator">
                                <div class="recording-dot"></div>
                                <span>Recording... <span id="recordingTime">00:00</span></span>
                            </div>

                            <div class="audio-preview" id="audioPreview">
                                <audio id="audioPlayback" controls></audio>
                                <div class="audio-preview-actions">
                                    <button type="button" class="btn-danger" onclick="discardRecording()">Discard</button>
                                    <button type="button" onclick="useRecording()">Use This Recording</button>
                                </div>
                            </div>
                        </div>
                    </div>
                    {{end}}

                    <!-- Thumbnail Upload -->
                    <div class="form-group" id="thumbnailUploadGroup" style="display: none;">
                        <label>Thumbnail / Cover Image (optional)</label>
                        <input type="file" name="thumbnail" id="thumbnail" accept="image/*" style="display: none;">
                        <div class="custom-file-upload" onclick="document.getElementById('thumbnail').click()">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 8px;">
                                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                                <circle cx="8.5" cy="8.5" r="1.5"></circle>
                                <polyline points="21 15 16 10 5 21"></polyline>
                            </svg>
                            <span id="thumbnailFileName">Choose thumbnail image</span>
                        </div>
                        <div class="preview-container" id="thumbnailPreview">
                            <img src="" alt="Thumbnail Preview" class="preview-image" id="thumbnailPreviewImg">
                        </div>
                        <div class="file-info">Recommended: Square image for best display</div>
                    </div>

                    <button type="submit" class="full-width">Create Post</button>
                </form>
            </div>
            {{end}}
        </main>
    </div>

    <!-- Edit Modal -->
    <div id="editModal" class="modal">
        <div class="modal-content">
            <div class="modal-header"><h2>Edit Post</h2></div>
            <form method="POST" action="/admin/update" enctype="multipart/form-data">
                <input type="hidden" name="id" id="editId">
                <div class="form-group">
                    <label for="editContent">Content (max 2000 characters)</label>
                    <textarea name="content" id="editContent" maxlength="2000" required></textarea>
                    <div class="file-info"><span id="editCharCount">0</span> / 2000 characters</div>
                </div>
                <div class="form-group">
                    <label>Media (optional)</label>
                    <div style="display: flex; gap: 20px; margin-top: 8px;">
                        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                            <input type="radio" name="media_type" value="photo" checked class="editMediaTypeRadio">
                            <span>Photo</span>
                        </label>
                        {{if .EnableAudioUploads}}
                        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                            <input type="radio" name="media_type" value="audio" class="editMediaTypeRadio">
                            <span>Audio</span>
                        </label>
                        {{end}}
                        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                            <input type="radio" name="media_type" value="video" class="editMediaTypeRadio">
                            <span>Video</span>
                        </label>
                    </div>
                </div>
                <div class="form-group">
                    <input type="file" name="media" id="editMedia" accept="image/*" style="display: none;">
                    <div class="custom-file-upload" onclick="document.getElementById('editMedia').click()">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 8px;">
                            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                            <polyline points="17 8 12 3 7 8"></polyline>
                            <line x1="12" y1="3" x2="12" y2="15"></line>
                        </svg>
                        <span id="editMediaFileName">Choose file</span>
                    </div>
                    <div class="preview-container" id="editPhotoPreview">
                        <img src="" alt="Preview" class="preview-image" id="editPreviewImg">
                    </div>
                </div>
                <div class="form-group" id="editThumbnailUploadGroup" style="display: none;">
                    <label>Thumbnail / Cover Image (optional)</label>
                    <input type="file" name="thumbnail" id="editThumbnail" accept="image/*" style="display: none;">
                    <div class="custom-file-upload" onclick="document.getElementById('editThumbnail').click()">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 8px;">
                            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                            <circle cx="8.5" cy="8.5" r="1.5"></circle>
                            <polyline points="21 15 16 10 5 21"></polyline>
                        </svg>
                        <span id="editThumbnailFileName">Choose thumbnail image</span>
                    </div>
                    <div class="preview-container" id="editThumbnailPreview">
                        <img src="" alt="Thumbnail Preview" class="preview-image" id="editThumbnailPreviewImg">
                    </div>
                    <div class="file-info">Recommended: Square image for best display</div>
                    <label style="display: flex; align-items: center; gap: 8px; margin-top: 8px; cursor: pointer;">
                        <input type="checkbox" name="remove_thumbnail" value="1" id="removeThumbnailCheckbox">
                        <span style="font-size: 13px; color: #8e8e8e;">Remove existing thumbnail</span>
                    </label>
                </div>
                <div class="modal-actions">
                    <button type="button" class="btn-secondary" onclick="closeEditModal()">Cancel</button>
                    <button type="submit">Update Post</button>
                </div>
            </form>
        </div>
    </div>

    <!-- Delete Modal -->
    <div id="deleteModal" class="modal">
        <div class="modal-content">
            <div class="modal-header"><h2>Delete Post</h2></div>
            <p style="margin-bottom: 20px;">Are you sure you want to delete this post? This action cannot be undone.</p>
            <form method="POST" action="/admin/delete">
                <input type="hidden" name="id" id="deleteId">
                <div class="modal-actions">
                    <button type="button" class="btn-secondary" onclick="closeDeleteModal()">Cancel</button>
                    <button type="submit" class="btn-danger">Delete Post</button>
                </div>
            </form>
        </div>
    </div>

    <script>
        // Mobile sidebar toggle
        function toggleSidebar() {
            document.getElementById('sidebar').classList.toggle('open');
            document.querySelector('.sidebar-overlay').classList.toggle('open');
        }

        // Auto-dismiss success messages
        document.addEventListener('DOMContentLoaded', function() {
            const messageDiv = document.querySelector('.message.success');
            if (messageDiv) {
                setTimeout(function() {
                    messageDiv.style.transition = 'opacity 0.5s ease-out';
                    messageDiv.style.opacity = '0';
                    setTimeout(function() {
                        messageDiv.remove();
                    }, 500);
                }, 3000);
            }
        });

        // Character counter
        const contentEl = document.getElementById('content');
        if (contentEl) {
            contentEl.addEventListener('input', function() {
                document.getElementById('charCount').textContent = this.value.length;
            });
        }

        const editContentEl = document.getElementById('editContent');
        if (editContentEl) {
            editContentEl.addEventListener('input', function() {
                document.getElementById('editCharCount').textContent = this.value.length;
            });
        }

        // Media type selector handler for create form
        document.querySelectorAll('.mediaTypeRadio').forEach(function(radio) {
            radio.addEventListener('change', function() {
                const mediaInput = document.getElementById('media');
                const type = this.value;
                const audioSourceOptions = document.getElementById('audioSourceOptions');
                const fileUploadGroup = document.getElementById('fileUploadGroup');
                const thumbnailUploadGroup = document.getElementById('thumbnailUploadGroup');

                if (type === 'photo') {
                    mediaInput.accept = 'image/*';
                    audioSourceOptions.classList.remove('active');
                    fileUploadGroup.style.display = 'block';
                    thumbnailUploadGroup.style.display = 'none';
                } else if (type === 'audio') {
                    mediaInput.accept = 'audio/mpeg,audio/mp3,audio/wav,audio/x-wav,audio/ogg,audio/aac,audio/mp4,audio/x-m4a,audio/webm,.mp3,.m4a,.wav,.ogg,.aac,.webm';
                    audioSourceOptions.classList.add('active');
                    document.querySelector('input[name="audio_source"][value="upload"]').checked = true;
                    document.getElementById('recordingUI').classList.remove('active');
                    fileUploadGroup.style.display = 'block';
                    thumbnailUploadGroup.style.display = 'block';
                } else if (type === 'video') {
                    mediaInput.accept = 'video/*';
                    audioSourceOptions.classList.remove('active');
                    fileUploadGroup.style.display = 'block';
                    thumbnailUploadGroup.style.display = 'block';
                }
            });
        });

        // Audio source selector handler
        document.querySelectorAll('.audioSourceRadio').forEach(function(radio) {
            radio.addEventListener('change', function() {
                const recordingUI = document.getElementById('recordingUI');
                const fileUploadGroup = document.getElementById('fileUploadGroup');

                if (this.value === 'record') {
                    recordingUI.classList.add('active');
                    fileUploadGroup.style.display = 'none';
                } else {
                    recordingUI.classList.remove('active');
                    fileUploadGroup.style.display = 'block';
                    discardRecording();
                }
            });
        });

        // Media type selector handler for edit form
        document.querySelectorAll('.editMediaTypeRadio').forEach(function(radio) {
            radio.addEventListener('change', function() {
                const mediaInput = document.getElementById('editMedia');
                const type = this.value;
                const editThumbnailUploadGroup = document.getElementById('editThumbnailUploadGroup');

                if (type === 'photo') {
                    mediaInput.accept = 'image/*';
                    editThumbnailUploadGroup.style.display = 'none';
                } else if (type === 'audio') {
                    mediaInput.accept = 'audio/mpeg,audio/mp3,audio/wav,audio/x-wav,audio/ogg,audio/aac,audio/mp4,audio/x-m4a,audio/webm,.mp3,.m4a,.wav,.ogg,.aac,.webm';
                    editThumbnailUploadGroup.style.display = 'block';
                } else if (type === 'video') {
                    mediaInput.accept = 'video/*';
                    editThumbnailUploadGroup.style.display = 'block';
                }
            });
        });

        // File change handlers
        const mediaEl = document.getElementById('media');
        if (mediaEl) {
            mediaEl.addEventListener('change', function(e) {
                const file = e.target.files[0];
                const fileNameSpan = document.getElementById('mediaFileName');

                if (file) {
                    fileNameSpan.textContent = file.name;
                    if (file.type.startsWith('image/')) {
                        const reader = new FileReader();
                        reader.onload = function(e) {
                            document.getElementById('previewImg').src = e.target.result;
                            document.getElementById('photoPreview').classList.add('active');
                        };
                        reader.readAsDataURL(file);
                    } else {
                        document.getElementById('photoPreview').classList.remove('active');
                    }
                } else {
                    fileNameSpan.textContent = 'Choose file';
                    document.getElementById('photoPreview').classList.remove('active');
                }
            });
        }

        const editMediaEl = document.getElementById('editMedia');
        if (editMediaEl) {
            editMediaEl.addEventListener('change', function(e) {
                const file = e.target.files[0];
                const fileNameSpan = document.getElementById('editMediaFileName');

                if (file) {
                    fileNameSpan.textContent = file.name;
                    if (file.type.startsWith('image/')) {
                        const reader = new FileReader();
                        reader.onload = function(e) {
                            document.getElementById('editPreviewImg').src = e.target.result;
                            document.getElementById('editPhotoPreview').classList.add('active');
                        };
                        reader.readAsDataURL(file);
                    } else {
                        document.getElementById('editPhotoPreview').classList.remove('active');
                    }
                } else {
                    fileNameSpan.textContent = 'Choose file';
                    document.getElementById('editPhotoPreview').classList.remove('active');
                }
            });
        }

        // Thumbnail file change handlers
        const thumbnailEl = document.getElementById('thumbnail');
        if (thumbnailEl) {
            thumbnailEl.addEventListener('change', function(e) {
                const file = e.target.files[0];
                const fileNameSpan = document.getElementById('thumbnailFileName');

                if (file) {
                    fileNameSpan.textContent = file.name;
                    if (file.type.startsWith('image/')) {
                        const reader = new FileReader();
                        reader.onload = function(e) {
                            document.getElementById('thumbnailPreviewImg').src = e.target.result;
                            document.getElementById('thumbnailPreview').classList.add('active');
                        };
                        reader.readAsDataURL(file);
                    }
                } else {
                    fileNameSpan.textContent = 'Choose thumbnail image';
                    document.getElementById('thumbnailPreview').classList.remove('active');
                }
            });
        }

        const editThumbnailEl = document.getElementById('editThumbnail');
        if (editThumbnailEl) {
            editThumbnailEl.addEventListener('change', function(e) {
                const file = e.target.files[0];
                const fileNameSpan = document.getElementById('editThumbnailFileName');

                if (file) {
                    fileNameSpan.textContent = file.name;
                    if (file.type.startsWith('image/')) {
                        const reader = new FileReader();
                        reader.onload = function(e) {
                            document.getElementById('editThumbnailPreviewImg').src = e.target.result;
                            document.getElementById('editThumbnailPreview').classList.add('active');
                        };
                        reader.readAsDataURL(file);
                    }
                } else {
                    fileNameSpan.textContent = 'Choose thumbnail image';
                    document.getElementById('editThumbnailPreview').classList.remove('active');
                }
            });
        }

        // Modal functions
        function openEditModal(id, content) {
            document.getElementById('editId').value = id;
            document.getElementById('editContent').value = content;
            document.getElementById('editCharCount').textContent = content.length;
            document.getElementById('editModal').classList.add('active');
        }

        function closeEditModal() {
            document.getElementById('editModal').classList.remove('active');
        }

        function openDeleteModal(id) {
            document.getElementById('deleteId').value = id;
            document.getElementById('deleteModal').classList.add('active');
        }

        function closeDeleteModal() {
            document.getElementById('deleteModal').classList.remove('active');
        }

        window.onclick = function(event) {
            if (event.target.classList.contains('modal')) {
                event.target.classList.remove('active');
            }
        }

        // Audio Recording functionality
        let mediaRecorder = null;
        let audioChunks = [];
        let recordedBlob = null;
        let recordingStartTime = null;
        let recordingTimer = null;

        async function startRecording() {
            try {
                const stream = await navigator.mediaDevices.getUserMedia({ audio: true });

                let mimeType = 'audio/webm';
                if (MediaRecorder.isTypeSupported('audio/webm;codecs=opus')) {
                    mimeType = 'audio/webm;codecs=opus';
                } else if (MediaRecorder.isTypeSupported('audio/webm')) {
                    mimeType = 'audio/webm';
                } else if (MediaRecorder.isTypeSupported('audio/ogg;codecs=opus')) {
                    mimeType = 'audio/ogg;codecs=opus';
                }

                mediaRecorder = new MediaRecorder(stream, { mimeType: mimeType });
                audioChunks = [];

                mediaRecorder.ondataavailable = function(event) {
                    if (event.data.size > 0) {
                        audioChunks.push(event.data);
                    }
                };

                mediaRecorder.onstop = function() {
                    recordedBlob = new Blob(audioChunks, { type: mimeType });
                    const audioUrl = URL.createObjectURL(recordedBlob);
                    document.getElementById('audioPlayback').src = audioUrl;
                    document.getElementById('audioPreview').classList.add('active');
                    document.getElementById('recordingIndicator').classList.remove('active');
                    stream.getTracks().forEach(track => track.stop());
                };

                mediaRecorder.start();
                recordingStartTime = Date.now();

                document.getElementById('startRecordBtn').style.display = 'none';
                document.getElementById('stopRecordBtn').style.display = 'inline-flex';
                document.getElementById('recordingIndicator').classList.add('active');
                document.getElementById('audioPreview').classList.remove('active');

                recordingTimer = setInterval(updateRecordingTime, 1000);
                updateRecordingTime();

            } catch (err) {
                console.error('Error accessing microphone:', err);
                alert('Could not access microphone. Please ensure you have granted microphone permissions.');
            }
        }

        function stopRecording() {
            if (mediaRecorder && mediaRecorder.state !== 'inactive') {
                mediaRecorder.stop();
            }

            document.getElementById('startRecordBtn').style.display = 'inline-flex';
            document.getElementById('stopRecordBtn').style.display = 'none';

            if (recordingTimer) {
                clearInterval(recordingTimer);
                recordingTimer = null;
            }
        }

        function updateRecordingTime() {
            if (recordingStartTime) {
                const elapsed = Math.floor((Date.now() - recordingStartTime) / 1000);
                const minutes = Math.floor(elapsed / 60).toString().padStart(2, '0');
                const seconds = (elapsed % 60).toString().padStart(2, '0');
                document.getElementById('recordingTime').textContent = minutes + ':' + seconds;
            }
        }

        function discardRecording() {
            recordedBlob = null;
            audioChunks = [];

            const audioPreviewEl = document.getElementById('audioPreview');
            const audioPlaybackEl = document.getElementById('audioPlayback');
            const recordingTimeEl = document.getElementById('recordingTime');

            if (audioPreviewEl) audioPreviewEl.classList.remove('active');
            if (audioPlaybackEl) audioPlaybackEl.src = '';
            if (recordingTimeEl) recordingTimeEl.textContent = '00:00';

            const mediaInput = document.getElementById('media');
            if (mediaInput) {
                mediaInput.value = '';
                document.getElementById('mediaFileName').textContent = 'Choose file';
            }
        }

        function useRecording() {
            if (!recordedBlob) {
                alert('No recording available');
                return;
            }

            const fileName = 'recording-' + Date.now() + '.webm';
            const file = new File([recordedBlob], fileName, { type: recordedBlob.type });

            const dataTransfer = new DataTransfer();
            dataTransfer.items.add(file);
            document.getElementById('media').files = dataTransfer.files;

            document.getElementById('mediaFileName').textContent = fileName;
            document.getElementById('audioPreview').innerHTML = '<div style="color: #155724; font-weight: 600;">✓ Recording ready to upload</div>';
        }

        // Form submit handler for audio recording
        const createFormEl = document.getElementById('createForm');
        if (createFormEl) {
            createFormEl.addEventListener('submit', function(e) {
                const audioSourceRadio = document.querySelector('input[name="audio_source"]:checked');
                const isRecordMode = audioSourceRadio && audioSourceRadio.value === 'record';
                const isRecording = mediaRecorder && mediaRecorder.state === 'recording';
                const hasRecordedBlob = recordedBlob !== null;
                const fileInput = document.getElementById('media');
                const hasFileSelected = fileInput.files && fileInput.files.length > 0;

                if (isRecordMode && isRecording) {
                    e.preventDefault();

                    const originalOnStop = mediaRecorder.onstop;
                    mediaRecorder.onstop = function(event) {
                        if (originalOnStop) {
                            originalOnStop.call(mediaRecorder, event);
                        }

                        setTimeout(function() {
                            if (recordedBlob) {
                                const fileName = 'recording-' + Date.now() + '.webm';
                                const file = new File([recordedBlob], fileName, { type: recordedBlob.type });
                                const dataTransfer = new DataTransfer();
                                dataTransfer.items.add(file);
                                document.getElementById('media').files = dataTransfer.files;
                            }
                            document.getElementById('createForm').submit();
                        }, 100);
                    };

                    mediaRecorder.stop();
                    if (recordingTimer) {
                        clearInterval(recordingTimer);
                        recordingTimer = null;
                    }
                    return;
                }

                if (isRecordMode && hasRecordedBlob && !hasFileSelected) {
                    e.preventDefault();
                    const fileName = 'recording-' + Date.now() + '.webm';
                    const file = new File([recordedBlob], fileName, { type: recordedBlob.type });
                    const dataTransfer = new DataTransfer();
                    dataTransfer.items.add(file);
                    document.getElementById('media').files = dataTransfer.files;
                    document.getElementById('createForm').submit();
                    return;
                }
            });
        }
    </script>
</body>
</html>`
