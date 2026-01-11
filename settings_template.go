package main

const settingsTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.PageTitle}} - Admin</title>
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
            margin-bottom: 20px;
        }

        /* Messages */
        .message { padding: 12px 16px; margin-bottom: 20px; border-radius: 8px; font-size: 14px; line-height: 18px; }
        .success { background-color: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .error { background-color: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }

        /* Section Title */
        .section-title {
            font-size: 18px;
            font-weight: 600;
            color: #262626;
            margin-bottom: 20px;
            padding-bottom: 12px;
            border-bottom: 1px solid #dbdbdb;
        }

        /* Forms */
        .form-group { margin-bottom: 20px; }
        label { display: block; margin-bottom: 8px; font-weight: 600; font-size: 14px; color: #262626; }
        input[type="text"],
        input[type="password"],
        select {
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
        input[type="password"]:hover,
        select:hover {
            border-color: #a8a8a8;
        }
        input[type="text"]:focus,
        input[type="password"]:focus,
        select:focus {
            outline: none;
            border-color: #0095f6;
            background-color: #ffffff;
        }
        .file-info { font-size: 12px; color: #8e8e8e; margin-top: 6px; }

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

        /* Settings Section */
        .settings-section {
            margin-bottom: 32px;
            padding-bottom: 32px;
            border-bottom: 1px solid #efefef;
        }
        .settings-section:last-child {
            border-bottom: none;
            margin-bottom: 0;
            padding-bottom: 0;
        }

        /* Avatar Upload */
        .avatar-upload-container {
            display: flex;
            align-items: flex-start;
            gap: 24px;
            margin-bottom: 24px;
        }
        .avatar-preview {
            width: 100px;
            height: 100px;
            border-radius: 50%;
            background-color: #000000;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 600;
            font-size: 40px;
            color: #ffffff;
            flex-shrink: 0;
            overflow: hidden;
        }
        .avatar-preview img {
            width: 100%;
            height: 100%;
            object-fit: cover;
        }
        .avatar-upload-info {
            flex: 1;
        }
        .custom-file-upload {
            display: inline-flex;
            align-items: center;
            padding: 10px 18px;
            background-color: #f0f0f0;
            color: #262626;
            border: 2px dashed #dbdbdb;
            border-radius: 8px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.3s ease;
        }
        .custom-file-upload:hover {
            background-color: #e8e8e8;
            border-color: #0095f6;
        }
        .custom-file-upload svg {
            margin-right: 8px;
        }

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
            .avatar-upload-container {
                flex-direction: column;
                align-items: center;
                text-align: center;
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
                        <a href="/admin">
                            <svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="16"></line><line x1="8" y1="12" x2="16" y2="12"></line></svg>
                            New Post
                        </a>
                    </li>
                    <li>
                        <a href="/admin?view=posts">
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
                        <a href="/admin/settings" class="{{if eq .View "site-info"}}active{{end}}">
                            <svg viewBox="0 0 24 24"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
                            Site Info
                        </a>
                    </li>
                    <li>
                        <a href="/admin/settings/appearance" class="{{if eq .View "appearance"}}active{{end}}">
                            <svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line></svg>
                            Appearance
                        </a>
                    </li>
                    <li>
                        <a href="/admin/settings/security" class="{{if eq .View "security"}}active{{end}}">
                            <svg viewBox="0 0 24 24"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path></svg>
                            Security
                        </a>
                    </li>
                    <li>
                        <a href="/admin/settings/backup" class="{{if eq .View "backup"}}active{{end}}">
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

            <div class="content-header">
                <h1>{{.PageTitle}}</h1>
            </div>

            {{if eq .View "site-info"}}
            <!-- Site Info Settings -->
            <div class="content-container">
                <form method="POST" action="/admin/settings/update" enctype="multipart/form-data" id="settingsForm">
                    <input type="hidden" name="section" value="site-info">

                    <!-- Avatar Upload -->
                    <div class="avatar-upload-container">
                        <div class="avatar-preview" id="avatarPreview">
                            {{if .Settings.AvatarPath}}
                                <img src="/uploads/{{.Settings.AvatarPath}}" alt="Avatar">
                            {{else}}
                                <div style="width: 100%; height: 100%; background: #ffffff; border-radius: 50%; display: flex; align-items: center; justify-content: center; color: #262626;">{{.Settings.UserInitial}}</div>
                            {{end}}
                        </div>
                        <div class="avatar-upload-info">
                            <label style="display: block; margin-bottom: 12px;">Avatar</label>
                            <input type="file" name="avatar" id="avatar" accept="image/jpeg,image/jpg,image/png" style="display: none;">
                            <div class="custom-file-upload" onclick="document.getElementById('avatar').click()">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                                    <polyline points="17 8 12 3 7 8"></polyline>
                                    <line x1="12" y1="3" x2="12" y2="15"></line>
                                </svg>
                                <span id="avatarFileName">Choose image</span>
                            </div>
                            <div class="file-info" style="margin-top: 8px;">JPG or PNG, will be resized to 200x200px</div>
                        </div>
                    </div>

                    {{if .Settings.AvatarPath}}
                    <div class="form-group">
                        <label>Display Preference</label>
                        <div style="display: flex; gap: 20px; margin-top: 8px;">
                            <label style="display: flex; align-items: center; gap: 8px; cursor: pointer; font-weight: normal;">
                                <input type="radio" name="avatar_preference" value="avatar" {{if eq .Settings.AvatarPreference "avatar"}}checked{{end}} style="cursor: pointer;">
                                <span>Show Avatar</span>
                            </label>
                            <label style="display: flex; align-items: center; gap: 8px; cursor: pointer; font-weight: normal;">
                                <input type="radio" name="avatar_preference" value="initials" {{if or (eq .Settings.AvatarPreference "initials") (eq .Settings.AvatarPreference "")}}checked{{end}} style="cursor: pointer;">
                                <span>Show Initials</span>
                            </label>
                        </div>
                        <div class="file-info">Choose whether to display your avatar or initials on posts</div>
                    </div>
                    {{end}}

                    <div class="form-group">
                        <label for="userInitial">Avatar Initial</label>
                        <input type="text" name="user_initial" id="userInitial" value="{{.Settings.UserInitial}}" maxlength="3" required>
                        <div class="file-info">Your initials displayed when showing initials (1-3 characters)</div>
                    </div>

                    <div class="form-group">
                        <label for="siteTitle">Site Title</label>
                        <input type="text" name="site_title" id="siteTitle" value="{{.Settings.SiteTitle}}" required>
                        <div class="file-info">The main title displayed on your blog</div>
                    </div>

                    <div class="form-group">
                        <label for="siteSubtitle">Site Subtitle</label>
                        <input type="text" name="site_subtitle" id="siteSubtitle" value="{{.Settings.SiteSubtitle}}" required>
                        <div class="file-info">Subtitle or tagline displayed below the title</div>
                    </div>

                    <button type="submit" class="full-width">Save Site Info</button>
                </form>
            </div>

            {{else if eq .View "appearance"}}
            <!-- Appearance Settings -->
            <div class="content-container">
                <form method="POST" action="/admin/settings/update" id="appearanceForm">
                    <input type="hidden" name="section" value="appearance">

                    <div class="form-group">
                        <label for="siteTheme">Theme</label>
                        <select name="site_theme" id="siteTheme" onchange="toggleCustomColors()">
                            <option value="default" {{if eq .Settings.SiteTheme "default"}}selected{{end}}>Light</option>
                            <option value="dark" {{if eq .Settings.SiteTheme "dark"}}selected{{end}}>Dark</option>
                            <option value="custom" {{if eq .Settings.SiteTheme "custom"}}selected{{end}}>Custom</option>
                        </select>
                        <div class="file-info">Choose the color scheme for your blog</div>
                    </div>

                    <!-- Custom Theme Colors -->
                    <div id="customColorsSection" style="display: {{if eq .Settings.SiteTheme "custom"}}block{{else}}none{{end}}; margin-top: 20px; padding: 20px; background: #f8f9fa; border-radius: 12px; border: 1px solid #e9ecef;">
                        <div class="section-title" style="margin-bottom: 16px;">Custom Theme Colors</div>

                        <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; margin-bottom: 20px;">
                            <div class="form-group" style="margin-bottom: 0;">
                                <label for="customBgColor" style="font-size: 13px;">Background</label>
                                <div style="display: flex; align-items: center; gap: 8px;">
                                    <input type="color" name="custom_bg_color" id="customBgColor" value="{{.Settings.CustomBgColor}}" onchange="updatePreview()" style="width: 50px; height: 36px; padding: 2px; border: 1px solid #dbdbdb; border-radius: 6px; cursor: pointer;">
                                    <input type="text" id="customBgColorText" value="{{.Settings.CustomBgColor}}" onchange="syncColorFromText('customBgColor')" style="flex: 1; font-size: 12px; font-family: monospace;">
                                </div>
                            </div>
                            <div class="form-group" style="margin-bottom: 0;">
                                <label for="customTextColor" style="font-size: 13px;">Text</label>
                                <div style="display: flex; align-items: center; gap: 8px;">
                                    <input type="color" name="custom_text_color" id="customTextColor" value="{{.Settings.CustomTextColor}}" onchange="updatePreview()" style="width: 50px; height: 36px; padding: 2px; border: 1px solid #dbdbdb; border-radius: 6px; cursor: pointer;">
                                    <input type="text" id="customTextColorText" value="{{.Settings.CustomTextColor}}" onchange="syncColorFromText('customTextColor')" style="flex: 1; font-size: 12px; font-family: monospace;">
                                </div>
                            </div>
                            <div class="form-group" style="margin-bottom: 0;">
                                <label for="customAccentColor" style="font-size: 13px;">Accent</label>
                                <div style="display: flex; align-items: center; gap: 8px;">
                                    <input type="color" name="custom_accent_color" id="customAccentColor" value="{{.Settings.CustomAccentColor}}" onchange="updatePreview()" style="width: 50px; height: 36px; padding: 2px; border: 1px solid #dbdbdb; border-radius: 6px; cursor: pointer;">
                                    <input type="text" id="customAccentColorText" value="{{.Settings.CustomAccentColor}}" onchange="syncColorFromText('customAccentColor')" style="flex: 1; font-size: 12px; font-family: monospace;">
                                </div>
                            </div>
                        </div>

                        <!-- Live Preview -->
                        <div style="margin-top: 16px;">
                            <label style="font-size: 13px; font-weight: 600; display: block; margin-bottom: 8px;">Preview</label>
                            <div id="themePreview" style="border-radius: 8px; overflow: hidden; border: 1px solid #dbdbdb;">
                                <div id="previewHeader" style="padding: 12px 16px; border-bottom: 1px solid;">
                                    <div style="display: flex; align-items: center; gap: 10px;">
                                        <div id="previewAvatar" style="width: 32px; height: 32px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: 600; font-size: 12px; color: white;">{{.Settings.UserInitial}}</div>
                                        <div>
                                            <div id="previewTitle" style="font-weight: 600; font-size: 14px;">{{.Settings.SiteTitle}}</div>
                                            <div id="previewSubtitle" style="font-size: 12px; opacity: 0.7;">{{.Settings.SiteSubtitle}}</div>
                                        </div>
                                    </div>
                                </div>
                                <div id="previewBody" style="padding: 16px;">
                                    <div id="previewCard" style="padding: 12px; border-radius: 8px; border: 1px solid;">
                                        <div id="previewPostTitle" style="font-weight: 600; margin-bottom: 4px;">Sample Post Title</div>
                                        <div id="previewPostText" style="font-size: 13px; opacity: 0.9;">This is how your blog content will appear with the selected colors.</div>
                                        <div id="previewLink" style="font-size: 12px; margin-top: 8px;">Read more →</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <button type="submit" class="full-width" style="margin-top: 20px;">Save Appearance</button>
                </form>
            </div>

            <script>
            function toggleCustomColors() {
                var theme = document.getElementById('siteTheme').value;
                var customSection = document.getElementById('customColorsSection');
                customSection.style.display = theme === 'custom' ? 'block' : 'none';
                if (theme === 'custom') {
                    updatePreview();
                }
            }

            function syncColorFromText(colorId) {
                var colorInput = document.getElementById(colorId);
                var textInput = document.getElementById(colorId + 'Text');
                var value = textInput.value.trim();
                if (/^#[0-9A-Fa-f]{6}$/.test(value)) {
                    colorInput.value = value;
                    updatePreview();
                }
            }

            function updatePreview() {
                var bg = document.getElementById('customBgColor').value;
                var text = document.getElementById('customTextColor').value;
                var accent = document.getElementById('customAccentColor').value;

                // Update text inputs
                document.getElementById('customBgColorText').value = bg;
                document.getElementById('customTextColorText').value = text;
                document.getElementById('customAccentColorText').value = accent;

                // Calculate derived colors
                var headerBg = adjustBrightness(bg, isLightColor(bg) ? -5 : 10);
                var borderColor = adjustBrightness(bg, isLightColor(bg) ? -15 : 20);
                var cardBg = adjustBrightness(bg, isLightColor(bg) ? 3 : 5);

                // Apply to preview
                document.getElementById('previewBody').style.backgroundColor = bg;
                document.getElementById('previewBody').style.color = text;
                document.getElementById('previewHeader').style.backgroundColor = headerBg;
                document.getElementById('previewHeader').style.color = text;
                document.getElementById('previewHeader').style.borderColor = borderColor;
                document.getElementById('previewCard').style.backgroundColor = cardBg;
                document.getElementById('previewCard').style.borderColor = borderColor;
                document.getElementById('previewAvatar').style.backgroundColor = accent;
                document.getElementById('previewLink').style.color = accent;
                document.getElementById('previewTitle').style.color = text;
                document.getElementById('previewSubtitle').style.color = text;
                document.getElementById('previewPostTitle').style.color = text;
                document.getElementById('previewPostText').style.color = text;
            }

            function isLightColor(hex) {
                var r = parseInt(hex.slice(1,3), 16);
                var g = parseInt(hex.slice(3,5), 16);
                var b = parseInt(hex.slice(5,7), 16);
                var luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
                return luminance > 0.5;
            }

            function adjustBrightness(hex, percent) {
                var r = parseInt(hex.slice(1,3), 16);
                var g = parseInt(hex.slice(3,5), 16);
                var b = parseInt(hex.slice(5,7), 16);
                r = Math.min(255, Math.max(0, r + Math.round(r * percent / 100)));
                g = Math.min(255, Math.max(0, g + Math.round(g * percent / 100)));
                b = Math.min(255, Math.max(0, b + Math.round(b * percent / 100)));
                return '#' + [r,g,b].map(x => x.toString(16).padStart(2,'0')).join('');
            }

            // Initialize preview on load
            if (document.getElementById('siteTheme').value === 'custom') {
                updatePreview();
            }
            </script>

            {{else if eq .View "security"}}
            <!-- Security Settings -->
            <div class="content-container">
                <form method="POST" action="/admin/settings/update" id="securityForm">
                    <input type="hidden" name="section" value="security">

                    <!-- Public Access -->
                    <div class="settings-section">
                        <div class="section-title">Public Access Control</div>
                        {{if .Settings.HasViewerPassword}}
                        <div style="padding: 12px 16px; background-color: #d4edda; border: 1px solid #c3e6cb; border-radius: 8px; margin-bottom: 16px;">
                            <span style="color: #155724; font-size: 14px;">Blog is currently password-protected</span>
                        </div>
                        {{end}}
                        <div class="form-group">
                            <label for="viewerPassword">{{if .Settings.HasViewerPassword}}New Access Password{{else}}Set Access Password{{end}}</label>
                            <input type="password" name="viewer_password" id="viewerPassword" placeholder="{{if .Settings.HasViewerPassword}}Leave blank to keep current password{{else}}Leave blank to keep blog public{{end}}">
                            <div class="file-info">Require a password for visitors to view your blog. Leave blank to keep your blog publicly accessible.</div>
                        </div>
                        {{if .Settings.HasViewerPassword}}
                        <div class="form-group">
                            <label style="display: flex; align-items: center; gap: 8px; font-weight: normal; cursor: pointer;">
                                <input type="checkbox" name="remove_viewer_password" value="true">
                                Make blog publicly accessible (remove password)
                            </label>
                        </div>
                        {{end}}
                    </div>

                    <!-- Admin Password -->
                    <div class="settings-section">
                        <div class="section-title">Admin Password</div>
                        <div class="form-group">
                            <label for="adminPassword">New Admin Password</label>
                            <input type="password" name="admin_password" id="adminPassword" placeholder="Leave blank to keep current password">
                            <div class="file-info">Password for accessing the admin area. Leave blank to keep current password.</div>
                        </div>
                        <div class="form-group">
                            <label for="adminPasswordConfirm">Confirm Admin Password</label>
                            <input type="password" name="admin_password_confirm" id="adminPasswordConfirm" placeholder="Leave blank to keep current password">
                        </div>
                    </div>

                    <button type="submit" class="full-width">Save Security Settings</button>
                </form>
            </div>

            <!-- Custom Domain Section -->
            {{if .CanEnableCustomDomain}}
            <div class="content-container">
                <div class="section-title">Custom Domain</div>
                <p style="font-size: 14px; color: #8e8e8e; margin-bottom: 16px;">
                    Enable custom domain support to connect your own domain to this blog.
                </p>
                <form method="POST" action="/admin/domain/enable">
                    <button type="submit">Enable Custom Domains</button>
                </form>
                <p style="font-size: 12px; color: #8e8e8e; margin-top: 12px;">
                    This will register your current hostname for custom domain configuration.
                </p>
            </div>
            {{else if .CustomDomainEnabled}}
            <div class="content-container">
                <div class="section-title">Custom Domain</div>

                <div style="margin-bottom: 24px; padding: 12px 16px; background-color: #e7f3ff; border: 1px solid #b6d4fe; border-radius: 8px;">
                    <span style="color: #084298; font-size: 14px;">Default hostname: <strong>{{.InstanceHostname}}</strong></span>
                </div>

                {{if .CustomDomain}}
                    {{if .CustomDomain.ActivatedAt.Valid}}
                    <div style="padding: 16px; background-color: #d4edda; border: 1px solid #c3e6cb; border-radius: 8px; margin-bottom: 20px;">
                        <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 12px;">
                            <span style="color: #155724; font-size: 18px;">&#10003;</span>
                            <span style="color: #155724; font-size: 16px; font-weight: 600;">Active</span>
                        </div>
                        <p style="color: #155724; font-size: 14px; margin-bottom: 8px;">
                            Domain: <strong>{{.CustomDomain.Domain}}</strong>
                        </p>
                        <p style="color: #155724; font-size: 13px;">
                            Your blog is accessible at:
                        </p>
                        <ul style="color: #155724; font-size: 13px; margin: 8px 0 0 20px;">
                            <li>https://{{.CustomDomain.Domain}} (custom domain)</li>
                            {{if .InstanceHostname}}<li>https://{{.InstanceHostname}} (default)</li>{{end}}
                        </ul>
                    </div>
                    <form method="POST" action="/admin/domain/remove" onsubmit="return confirm('Are you sure you want to remove this custom domain?');">
                        <button type="submit" class="btn-danger">Remove Domain</button>
                    </form>

                    {{else if .CustomDomain.VerifiedAt.Valid}}
                    <div style="padding: 16px; background-color: #d1ecf1; border: 1px solid #bee5eb; border-radius: 8px; margin-bottom: 20px;">
                        <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 12px;">
                            <span style="color: #0c5460; font-size: 18px;">&#10003;</span>
                            <span style="color: #0c5460; font-size: 16px; font-weight: 600;">Verified - Activation Pending</span>
                        </div>
                        <p style="color: #0c5460; font-size: 14px;">
                            Domain: <strong>{{.CustomDomain.Domain}}</strong>
                        </p>
                    </div>
                    <div style="display: flex; gap: 12px;">
                        <form method="POST" action="/admin/domain/activate">
                            <button type="submit">Activate Domain</button>
                        </form>
                        <form method="POST" action="/admin/domain/remove" onsubmit="return confirm('Are you sure you want to remove this domain?');">
                            <button type="submit" class="btn-secondary">Cancel</button>
                        </form>
                    </div>

                    {{else}}
                    <div style="padding: 16px; background-color: #fff3cd; border: 1px solid #ffecb5; border-radius: 8px; margin-bottom: 20px;">
                        <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 12px;">
                            <span style="color: #664d03; font-size: 18px;">&#8987;</span>
                            <span style="color: #664d03; font-size: 16px; font-weight: 600;">Pending Verification</span>
                        </div>
                        <p style="color: #664d03; font-size: 14px; margin-bottom: 16px;">
                            Domain: <strong>{{.CustomDomain.Domain}}</strong>
                        </p>
                    </div>

                    <div style="margin-bottom: 24px;">
                        <h4 style="font-size: 15px; font-weight: 600; margin-bottom: 12px;">Add these DNS records at your domain provider:</h4>

                        <div style="margin-bottom: 16px; padding: 16px; background-color: #f8f9fa; border: 1px solid #e9ecef; border-radius: 8px;">
                            <p style="font-size: 13px; font-weight: 600; margin-bottom: 8px; color: #495057;">Step 1: Verification Record (TXT)</p>
                            <table style="width: 100%; font-size: 13px; border-collapse: collapse;">
                                <tr>
                                    <td style="padding: 4px 0; color: #6c757d; width: 60px;">Type:</td>
                                    <td style="padding: 4px 0; font-family: monospace;">TXT</td>
                                </tr>
                                <tr>
                                    <td style="padding: 4px 0; color: #6c757d;">Name:</td>
                                    <td style="padding: 4px 0; font-family: monospace; word-break: break-all;">_postastiq-verify.{{.CustomDomain.Domain}}</td>
                                </tr>
                                <tr>
                                    <td style="padding: 4px 0; color: #6c757d;">Value:</td>
                                    <td style="padding: 4px 0; font-family: monospace; word-break: break-all;">{{.CustomDomain.VerificationToken}}</td>
                                </tr>
                            </table>
                        </div>

                        <div style="margin-bottom: 16px; padding: 16px; background-color: #f8f9fa; border: 1px solid #e9ecef; border-radius: 8px;">
                            <p style="font-size: 13px; font-weight: 600; margin-bottom: 8px; color: #495057;">Step 2: Routing Record (CNAME)</p>
                            <table style="width: 100%; font-size: 13px; border-collapse: collapse;">
                                <tr>
                                    <td style="padding: 4px 0; color: #6c757d; width: 60px;">Type:</td>
                                    <td style="padding: 4px 0; font-family: monospace;">CNAME</td>
                                </tr>
                                <tr>
                                    <td style="padding: 4px 0; color: #6c757d;">Name:</td>
                                    <td style="padding: 4px 0; font-family: monospace;">{{.CustomDomain.Domain}}</td>
                                </tr>
                                <tr>
                                    <td style="padding: 4px 0; color: #6c757d;">Value:</td>
                                    <td style="padding: 4px 0; font-family: monospace;">{{.InstanceHostname}}</td>
                                </tr>
                            </table>
                        </div>

                        <p style="font-size: 13px; color: #6c757d; margin-bottom: 16px;">
                            DNS changes can take up to 48 hours to propagate.
                        </p>
                    </div>

                    <div style="display: flex; gap: 12px; flex-wrap: wrap;">
                        <form method="POST" action="/admin/domain/verify">
                            <button type="submit">Verify Domain</button>
                        </form>
                        <form method="POST" action="/admin/domain/remove" onsubmit="return confirm('Are you sure you want to cancel and remove this domain?');">
                            <button type="submit" class="btn-secondary">Cancel</button>
                        </form>
                    </div>
                    {{if gt .CustomDomain.VerificationAttempts 0}}
                    <p style="font-size: 12px; color: #6c757d; margin-top: 12px;">
                        Verification attempts: {{.CustomDomain.VerificationAttempts}}/5 this hour
                    </p>
                    {{end}}
                    {{end}}

                {{else}}
                    <p style="font-size: 14px; color: #8e8e8e; margin-bottom: 16px;">
                        Connect your own domain to this blog. You'll need access to your domain's DNS settings.
                    </p>
                    <form method="POST" action="/admin/domain/add">
                        <div class="form-group">
                            <label for="domain">Domain Name</label>
                            <input type="text" name="domain" id="domain" placeholder="blog.example.com" required>
                            <div class="file-info">Enter the domain you want to connect (e.g., blog.yourdomain.com)</div>
                        </div>
                        <button type="submit">Add Domain</button>
                    </form>
                {{end}}
            </div>
            {{end}}

            {{else if eq .View "backup"}}
            <!-- Backup Settings -->
            <div class="content-container">
                <!-- Download Backup -->
                <div class="settings-section">
                    <div class="section-title">Download Backup</div>
                    <p style="font-size: 14px; color: #8e8e8e; margin-bottom: 16px;">
                        Download a complete backup of your blog including all posts, settings, and media files.
                    </p>
                    <a href="/admin/backup" class="btn">Download Backup</a>
                    <div class="file-info" style="margin-top: 12px;">
                        Creates a ZIP file containing your database and all uploaded media files.
                    </div>
                </div>

                <!-- Restore -->
                <div>
                    <div class="section-title">Restore from Backup</div>
                    <p style="font-size: 14px; color: #8e8e8e; margin-bottom: 16px;">
                        Upload a backup ZIP file to restore your blog. This will replace all current data.
                    </p>
                    <form method="POST" action="/admin/restore" enctype="multipart/form-data" onsubmit="return confirmRestore()">
                        <div style="display: flex; gap: 12px; align-items: center; flex-wrap: wrap;">
                            <input type="file" name="backup_file" id="restoreFile" accept=".zip" required style="display: none;">
                            <div class="custom-file-upload" onclick="document.getElementById('restoreFile').click()">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                                    <polyline points="17 8 12 3 7 8"></polyline>
                                    <line x1="12" y1="3" x2="12" y2="15"></line>
                                </svg>
                                <span id="restoreFileName">Choose backup file</span>
                            </div>
                            <button type="submit" class="btn-danger">Restore Backup</button>
                        </div>
                    </form>
                    <div class="file-info" style="margin-top: 12px; color: #ed4956;">
                        Warning: Restoring will replace all current posts, settings, and media files.
                    </div>
                </div>
            </div>
            {{end}}
        </main>
    </div>

    <script>
        // Mobile sidebar toggle
        function toggleSidebar() {
            document.getElementById('sidebar').classList.toggle('open');
            document.querySelector('.sidebar-overlay').classList.toggle('open');
        }

        function confirmRestore() {
            return confirm('Are you sure you want to restore from this backup? This will replace ALL current data including posts, settings, and media files. This action cannot be undone.');
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

            // Restore file name display
            const restoreFile = document.getElementById('restoreFile');
            if (restoreFile) {
                restoreFile.addEventListener('change', function(e) {
                    const fileName = e.target.files[0] ? e.target.files[0].name : 'Choose backup file';
                    document.getElementById('restoreFileName').textContent = fileName;
                });
            }

            // Security form password validation
            const securityForm = document.getElementById('securityForm');
            if (securityForm) {
                securityForm.addEventListener('submit', function(e) {
                    const adminPass = document.getElementById('adminPassword').value;
                    const adminPassConfirm = document.getElementById('adminPasswordConfirm').value;

                    if (adminPass || adminPassConfirm) {
                        if (adminPass !== adminPassConfirm) {
                            e.preventDefault();
                            alert('Admin passwords do not match');
                            return false;
                        }
                    }
                });
            }

            // Avatar preview
            const avatar = document.getElementById('avatar');
            if (avatar) {
                avatar.addEventListener('change', function(e) {
                    const file = e.target.files[0];
                    const fileNameSpan = document.getElementById('avatarFileName');
                    const avatarPreview = document.getElementById('avatarPreview');

                    if (file) {
                        fileNameSpan.textContent = file.name;

                        const reader = new FileReader();
                        reader.onload = function(e) {
                            avatarPreview.innerHTML = '<img src="' + e.target.result + '" alt="Avatar preview">';
                        };
                        reader.readAsDataURL(file);

                        // Auto-submit
                        document.getElementById('settingsForm').submit();
                    } else {
                        fileNameSpan.textContent = 'Choose image';
                    }
                });
            }
        });
    </script>
</body>
</html>`
