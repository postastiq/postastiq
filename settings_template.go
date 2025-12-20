package main

const settingsTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Settings - Admin</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background-color: #fafafa;
            color: #262626;
            line-height: 1.6;
            min-height: 100vh;
            padding: 20px;
        }
        .container {
            max-width: 900px;
            width: 100%;
            margin: 0 auto;
            background-color: #ffffff;
            border: 1px solid #dbdbdb;
            border-radius: 8px;
            padding: 32px;
            margin-bottom: 20px;
        }
        h1 { font-size: 28px; font-weight: 600; color: #262626; margin-bottom: 8px; text-align: center; }
        .subtitle { font-size: 14px; color: #8e8e8e; text-align: center; margin-bottom: 32px; }
        .nav-links {
            display: flex;
            gap: 16px;
            justify-content: center;
            margin-bottom: 32px;
            padding-bottom: 24px;
            border-bottom: 1px solid #dbdbdb;
        }
        .nav-link {
            color: #0095f6;
            text-decoration: none;
            font-size: 14px;
            font-weight: 600;
            padding: 8px 16px;
            border-radius: 8px;
            transition: background-color 0.2s;
        }
        .nav-link:hover {
            background-color: #f0f0f0;
        }
        .message { padding: 12px 16px; margin-bottom: 20px; border-radius: 8px; font-size: 14px; line-height: 18px; }
        .success { background-color: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .error { background-color: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
        .section-title {
            font-size: 18px;
            font-weight: 600;
            color: #262626;
            margin-bottom: 20px;
            padding-bottom: 12px;
            border-bottom: 1px solid #dbdbdb;
        }
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
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
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
        button {
            background-color: #000000;
            color: #ffffff;
            padding: 12px 24px;
            border: none;
            cursor: pointer;
            border-radius: 8px;
            font-size: 14px;
            font-weight: 600;
            transition: transform 0.2s;
        }
        button:hover { transform: translateY(-1px); }
        button.full-width { width: 100%; }
        .settings-section {
            margin-bottom: 32px;
            padding-bottom: 32px;
            border-bottom: 1px solid #efefef;
        }
        .settings-section:last-child {
            border-bottom: none;
        }
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
        @media (max-width: 768px) {
            body { padding: 10px; }
            .container { padding: 20px; }
            .avatar-upload-container {
                flex-direction: column;
                align-items: center;
                text-align: center;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Settings</h1>
        <p class="subtitle">Manage site settings and configuration</p>

        <div class="nav-links">
            <a href="/admin" class="nav-link">Back to Admin</a>
            <a href="/rss" class="nav-link" target="_blank">RSS</a>
            <a href="/" class="nav-link" target="_blank">View Site</a>
        </div>

        {{if .Message}}
        <div class="message {{.MessageType}}">{{.Message}}</div>
        {{end}}

        <form method="POST" action="/admin/settings/update" enctype="multipart/form-data">
            <!-- Site Information -->
            <div class="settings-section">
                <div class="section-title">Site Information</div>

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

                <!-- Avatar Preference Radio Buttons - Only show if avatar is uploaded -->
                {{if .Settings.AvatarPath}}
                <div class="form-group">
                    <label>Display Preference</label>
                    <div style="display: flex; gap: 20px; margin-top: 8px;">
                        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                            <input type="radio" name="avatar_preference" value="avatar" {{if eq .Settings.AvatarPreference "avatar"}}checked{{end}} style="cursor: pointer;">
                            <span>Show Avatar</span>
                        </label>
                        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
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
            </div>

            <!-- Theme -->
            <div class="settings-section">
                <div class="section-title">Appearance</div>
                <div class="form-group">
                    <label for="siteTheme">Theme</label>
                    <select name="site_theme" id="siteTheme">
                        <option value="default" {{if eq .Settings.SiteTheme "default"}}selected{{end}}>Light Theme (Default)</option>
                        <option value="dark" {{if eq .Settings.SiteTheme "dark"}}selected{{end}}>Dark Theme</option>
                    </select>
                    <div class="file-info">Choose the color scheme for your blog</div>
                </div>
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

            <!-- Public Access -->
            <div class="settings-section">
                <div class="section-title">Public Access Control</div>
                {{if .Settings.HasViewerPassword}}
                <div style="padding: 12px 16px; background-color: #d4edda; border: 1px solid #c3e6cb; border-radius: 8px; margin-bottom: 16px;">
                    <span style="color: #155724; font-size: 14px;">✓ Blog is currently password-protected</span>
                </div>
                {{end}}
                <div class="form-group">
                    <label for="viewerPassword">{{if .Settings.HasViewerPassword}}New Access Password{{else}}Set Access Password{{end}}</label>
                    <input type="password" name="viewer_password" id="viewerPassword" placeholder="{{if .Settings.HasViewerPassword}}Leave blank to keep current password{{else}}Leave blank to keep blog public{{end}}">
                    <div class="file-info">Require a password for visitors to view your blog. Leave blank to keep your blog publicly accessible.</div>
                </div>
                {{if .Settings.HasViewerPassword}}
                <div class="form-group">
                    <label>
                        <input type="checkbox" name="remove_viewer_password" value="true">
                        Make blog publicly accessible (remove password)
                    </label>
                </div>
                {{end}}
            </div>

            <button type="submit" class="full-width">Save Settings</button>
        </form>
    </div>

    <!-- Custom Domain Section - Only shown for *.postastiq.com subdomains -->
    {{if .CanEnableCustomDomain}}
    <!-- Show Enable Button when not yet enabled but on postastiq.com subdomain -->
    <div class="container">
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
    <div class="container">
        <div class="section-title">Custom Domain</div>

        <div style="margin-bottom: 24px; padding: 12px 16px; background-color: #e7f3ff; border: 1px solid #b6d4fe; border-radius: 8px;">
            <span style="color: #084298; font-size: 14px;">Default hostname: <strong>{{.InstanceHostname}}</strong></span>
        </div>

        {{if .CustomDomain}}
            {{if .CustomDomain.ActivatedAt.Valid}}
            <!-- Domain Active State -->
            <div style="padding: 16px; background-color: #d4edda; border: 1px solid #c3e6cb; border-radius: 8px; margin-bottom: 20px;">
                <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 12px;">
                    <span style="color: #155724; font-size: 18px;">✓</span>
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
                <button type="submit" style="background-color: #dc3545;">Remove Domain</button>
            </form>

            {{else if .CustomDomain.VerifiedAt.Valid}}
            <!-- Domain Verified but Not Activated -->
            <div style="padding: 16px; background-color: #d1ecf1; border: 1px solid #bee5eb; border-radius: 8px; margin-bottom: 20px;">
                <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 12px;">
                    <span style="color: #0c5460; font-size: 18px;">✓</span>
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
                    <button type="submit" style="background-color: #6c757d;">Cancel</button>
                </form>
            </div>

            {{else}}
            <!-- Domain Pending Verification -->
            <div style="padding: 16px; background-color: #fff3cd; border: 1px solid #ffecb5; border-radius: 8px; margin-bottom: 20px;">
                <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 12px;">
                    <span style="color: #664d03; font-size: 18px;">⏳</span>
                    <span style="color: #664d03; font-size: 16px; font-weight: 600;">Pending Verification</span>
                </div>
                <p style="color: #664d03; font-size: 14px; margin-bottom: 16px;">
                    Domain: <strong>{{.CustomDomain.Domain}}</strong>
                </p>
            </div>

            <div style="margin-bottom: 24px;">
                <h4 style="font-size: 15px; font-weight: 600; margin-bottom: 12px;">Add these DNS records at your domain provider:</h4>

                <!-- TXT Record -->
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

                <!-- CNAME Record -->
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
                    ⚠️ DNS changes can take up to 48 hours to propagate.
                </p>
            </div>

            <div style="display: flex; gap: 12px; flex-wrap: wrap;">
                <form method="POST" action="/admin/domain/verify">
                    <button type="submit">Verify Domain</button>
                </form>
                <form method="POST" action="/admin/domain/remove" onsubmit="return confirm('Are you sure you want to cancel and remove this domain?');">
                    <button type="submit" style="background-color: #6c757d;">Cancel</button>
                </form>
            </div>
            {{if gt .CustomDomain.VerificationAttempts 0}}
            <p style="font-size: 12px; color: #6c757d; margin-top: 12px;">
                Verification attempts: {{.CustomDomain.VerificationAttempts}}/5 this hour
            </p>
            {{end}}
            {{end}}

        {{else}}
            <!-- No Domain Configured - Show Add Form -->
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

    <!-- Backup & Restore Section (outside main form) -->
    <div class="container">
        <div class="section-title">Backup & Restore</div>

        <!-- Backup -->
        <div style="margin-bottom: 32px; padding-bottom: 32px; border-bottom: 1px solid #efefef;">
            <h3 style="font-size: 16px; font-weight: 600; margin-bottom: 12px;">Download Backup</h3>
            <p style="font-size: 14px; color: #8e8e8e; margin-bottom: 16px;">
                Download a complete backup of your blog including all posts, settings, and media files.
            </p>
            <a href="/admin/backup" style="display: inline-block; background-color: #000000; color: #ffffff; padding: 12px 24px; text-decoration: none; border-radius: 8px; font-size: 14px; font-weight: 600; transition: transform 0.2s;">
                Download Backup
            </a>
            <div class="file-info" style="margin-top: 12px;">
                Creates a ZIP file containing your database and all uploaded media files.
            </div>
        </div>

        <!-- Restore -->
        <div>
            <h3 style="font-size: 16px; font-weight: 600; margin-bottom: 12px;">Restore from Backup</h3>
            <p style="font-size: 14px; color: #8e8e8e; margin-bottom: 16px;">
                Upload a backup ZIP file to restore your blog. This will replace all current data.
            </p>
            <form method="POST" action="/admin/restore" enctype="multipart/form-data" onsubmit="return confirmRestore()">
                <div style="display: flex; gap: 12px; align-items: center; flex-wrap: wrap;">
                    <input type="file" name="backup_file" id="restoreFile" accept=".zip" required style="display: none;">
                    <div class="custom-file-upload" onclick="document.getElementById('restoreFile').click()" style="display: inline-flex; align-items: center; padding: 10px 18px; background-color: #f0f0f0; color: #262626; border: 2px dashed #dbdbdb; border-radius: 8px; cursor: pointer; font-size: 14px; font-weight: 500;">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 8px;">
                            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                            <polyline points="17 8 12 3 7 8"></polyline>
                            <line x1="12" y1="3" x2="12" y2="15"></line>
                        </svg>
                        <span id="restoreFileName">Choose backup file</span>
                    </div>
                    <button type="submit" style="background-color: #dc3545; color: #ffffff; padding: 12px 24px; border: none; cursor: pointer; border-radius: 8px; font-size: 14px; font-weight: 600;">
                        Restore Backup
                    </button>
                </div>
            </form>
            <div class="file-info" style="margin-top: 12px; color: #dc3545;">
                ⚠️ Warning: Restoring will replace all current posts, settings, and media files.
            </div>
        </div>
    </div>

    <script>
        function confirmRestore() {
            return confirm('Are you sure you want to restore from this backup? This will replace ALL current data including posts, settings, and media files. This action cannot be undone.');
        }

        document.getElementById('restoreFile').addEventListener('change', function(e) {
            const fileName = e.target.files[0] ? e.target.files[0].name : 'Choose backup file';
            document.getElementById('restoreFileName').textContent = fileName;
        });

        // Auto-dismiss success messages after 3 seconds
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

        // Password confirmation validation
        document.querySelector('form').addEventListener('submit', function(e) {
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

        // Avatar preview and auto-save
        document.getElementById('avatar').addEventListener('change', function(e) {
            const file = e.target.files[0];
            const fileNameSpan = document.getElementById('avatarFileName');
            const avatarPreview = document.getElementById('avatarPreview');

            if (file) {
                fileNameSpan.textContent = file.name;

                // Show preview
                const reader = new FileReader();
                reader.onload = function(e) {
                    avatarPreview.innerHTML = '<img src="' + e.target.result + '" alt="Avatar preview">';
                };
                reader.readAsDataURL(file);

                // Auto-submit the form to save the avatar
                document.querySelector('form').submit();
            } else {
                fileNameSpan.textContent = 'Choose image';
            }
        });
    </script>
</body>
</html>`
