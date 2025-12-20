package main

const adminTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin</title>
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
        .message { padding: 12px 16px; margin-bottom: 20px; border-radius: 8px; font-size: 14px; line-height: 18px; }
        .success { background-color: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .error { background-color: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
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
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
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
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            resize: vertical;
            min-height: 100px;
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
        button.btn-danger { background-color: #ed4956; }
        button.btn-secondary { background-color: #8e8e8e; margin-left: 10px; }
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
            margin-bottom: 12px;
            padding-bottom: 12px;
            border-bottom: 1px solid #efefef;
        }
        .entry-id { font-weight: 600; font-size: 14px; }
        .entry-date { font-size: 12px; color: #8e8e8e; }
        .entry-content { margin-bottom: 12px; font-size: 14px; white-space: pre-wrap; }
        .entry-photo {
            max-width: 200px;
            max-height: 200px;
            border-radius: 4px;
            margin-bottom: 12px;
            cursor: pointer;
        }
        .entry-actions { display: flex; gap: 8px; }
        .entry-actions button { flex: 1; padding: 8px 16px; font-size: 13px; }
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
        .section-title {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 16px;
            padding-bottom: 8px;
            border-bottom: 2px solid #000000;
            display: flex;
            justify-content: space-between;
            align-items: center;
            cursor: pointer;
            user-select: none;
        }
        .toggle-icon {
            font-size: 20px;
            transition: transform 0.3s;
        }
        .toggle-icon.collapsed {
            transform: rotate(-90deg);
        }
        .collapsible-content {
            max-height: 2000px;
            overflow: hidden;
            transition: max-height 0.3s ease-out, opacity 0.3s ease-out;
            opacity: 1;
        }
        .collapsible-content.collapsed {
            max-height: 0;
            opacity: 0;
        }
        .empty-state { text-align: center; padding: 40px 20px; color: #8e8e8e; }
        .loading-indicator { text-align: center; padding: 20px; color: #8e8e8e; display: none; }
        .loading-indicator.active { display: block; }

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
        @media (max-width: 768px) {
            body { padding: 10px; }
            .container { padding: 20px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Admin</h1>
        <p class="subtitle">Create, edit, and publish with ease</p>

        <div style="display: flex; gap: 16px; justify-content: center; margin-bottom: 32px; padding-bottom: 24px; border-bottom: 1px solid #dbdbdb;">
            <a href="/admin/settings" style="color: #0095f6; text-decoration: none; font-size: 14px; font-weight: 600; padding: 8px 16px; border-radius: 8px; transition: background-color 0.2s;" onmouseover="this.style.backgroundColor='#f0f0f0'" onmouseout="this.style.backgroundColor='transparent'">Settings</a>
            <a href="/rss" target="_blank" style="color: #0095f6; text-decoration: none; font-size: 14px; font-weight: 600; padding: 8px 16px; border-radius: 8px; transition: background-color 0.2s;" onmouseover="this.style.backgroundColor='#f0f0f0'" onmouseout="this.style.backgroundColor='transparent'">RSS</a>
            <a href="/" target="_blank" style="color: #0095f6; text-decoration: none; font-size: 14px; font-weight: 600; padding: 8px 16px; border-radius: 8px; transition: background-color 0.2s;" onmouseover="this.style.backgroundColor='#f0f0f0'" onmouseout="this.style.backgroundColor='transparent'">View Site</a>
        </div>

        {{if .Message}}
        <div class="message {{.MessageType}}">{{.Message}}</div>
        {{end}}

        <div class="section-title" style="cursor: default;">Create New Entry</div>

        <form method="POST" action="/admin/create" enctype="multipart/form-data" id="createForm">
            <div class="form-group">
                <label for="content">Content (max 2000 characters)</label>
                <textarea name="content" id="content" maxlength="2000" required></textarea>
                <div class="file-info"><span id="charCount">0</span> / 2000 characters</div>
            </div>

            <div class="form-group">
                <label>Upload (optional)</label>
                <div style="display: flex; gap: 20px; margin-top: 8px;">
                    <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                        <input type="radio" name="media_type" value="photo" checked class="mediaTypeRadio">
                        <span>Photo</span>
                    </label>
                    <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                        <input type="radio" name="media_type" value="audio" class="mediaTypeRadio">
                        <span>Audio</span>
                    </label>
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

            <!-- Audio Source Options (shown when Audio is selected) -->
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

                <!-- Recording UI -->
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

            <!-- Thumbnail Upload (shown when Video or Audio is selected) -->
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

            <button type="submit" class="full-width">Create Entry</button>
        </form>
    </div>

    <div class="container">
        <div class="section-title" onclick="toggleAllEntries()">
            <span>All Entries</span>
            <span class="toggle-icon collapsed" id="allEntriesToggleIcon">▼</span>
        </div>
        <div class="collapsible-content collapsed" id="allEntriesContent">
            <div id="entriesContainer"></div>
            <div id="loadingIndicator" class="loading-indicator">Loading more entries...</div>
            <div id="emptyState" class="empty-state" style="display: none;">
                No entries found. Create your first entry above!
            </div>
        </div>
    </div>

    <div id="editModal" class="modal">
        <div class="modal-content">
            <div class="modal-header"><h2>Edit Entry</h2></div>
            <form method="POST" action="/admin/update" enctype="multipart/form-data">
                <input type="hidden" name="id" id="editId">
                <div class="form-group">
                    <label for="editContent">Content (max 2000 characters)</label>
                    <textarea name="content" id="editContent" maxlength="2000" required></textarea>
                    <div class="file-info"><span id="editCharCount">0</span> / 2000 characters</div>
                </div>
                <div class="form-group">
                    <label>Upload (optional)</label>
                    <div style="display: flex; gap: 20px; margin-top: 8px;">
                        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                            <input type="radio" name="media_type" value="photo" checked class="editMediaTypeRadio">
                            <span>Photo</span>
                        </label>
                        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                            <input type="radio" name="media_type" value="audio" class="editMediaTypeRadio">
                            <span>Audio</span>
                        </label>
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
                <!-- Thumbnail Upload for Edit (shown when Video or Audio is selected) -->
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
                    <button type="submit">Update Entry</button>
                </div>
            </form>
        </div>
    </div>

    <div id="deleteModal" class="modal">
        <div class="modal-content">
            <div class="modal-header"><h2>Delete Entry</h2></div>
            <p style="margin-bottom: 20px;">Are you sure you want to delete this entry?</p>
            <form method="POST" action="/admin/delete">
                <input type="hidden" name="id" id="deleteId">
                <div class="modal-actions">
                    <button type="button" class="btn-secondary" onclick="closeDeleteModal()">Cancel</button>
                    <button type="submit" class="btn-danger">Delete Entry</button>
                </div>
            </form>
        </div>
    </div>

    <script>
        let currentOffset = 0;
        let isLoading = false;
        let hasMore = true;

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

        function toggleAllEntries() {
            const content = document.getElementById('allEntriesContent');
            const icon = document.getElementById('allEntriesToggleIcon');

            content.classList.toggle('collapsed');
            icon.classList.toggle('collapsed');

            // Save the toggle state to localStorage
            const isCollapsed = content.classList.contains('collapsed');
            localStorage.setItem('allEntriesCollapsed', isCollapsed);
        }

        // Restore toggle state on page load
        function restoreToggleState() {
            const allEntriesCollapsed = localStorage.getItem('allEntriesCollapsed');
            if (allEntriesCollapsed === 'false') {
                const content = document.getElementById('allEntriesContent');
                const icon = document.getElementById('allEntriesToggleIcon');
                content.classList.remove('collapsed');
                icon.classList.remove('collapsed');
            }
        }

        // Call restore function on page load
        restoreToggleState();

        document.getElementById('content').addEventListener('input', function() {
            document.getElementById('charCount').textContent = this.value.length;
        });

        document.getElementById('editContent').addEventListener('input', function() {
            document.getElementById('editCharCount').textContent = this.value.length;
        });

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
                    // Reset to upload mode by default
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

        // Audio source selector handler (upload vs record)
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

        document.getElementById('media').addEventListener('change', function(e) {
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

        document.getElementById('editMedia').addEventListener('change', function(e) {
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

        // Thumbnail file change handler for create form
        document.getElementById('thumbnail').addEventListener('change', function(e) {
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

        // Thumbnail file change handler for edit form
        document.getElementById('editThumbnail').addEventListener('change', function(e) {
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

        async function loadEntries() {
            if (isLoading || !hasMore) return;

            isLoading = true;
            document.getElementById('loadingIndicator').classList.add('active');

            try {
                const response = await fetch('/admin/entries?offset=' + currentOffset);
                const data = await response.json();

                if (data.entries && data.entries.length > 0) {
                    const container = document.getElementById('entriesContainer');
                    data.entries.forEach(entry => {
                        const entryCard = createEntryCard(entry);
                        container.appendChild(entryCard);
                    });

                    currentOffset += data.entries.length;
                    hasMore = data.hasMore;
                } else {
                    hasMore = false;
                    if (currentOffset === 0) {
                        document.getElementById('emptyState').style.display = 'block';
                    }
                }
            } catch (error) {
                console.error('Error loading entries:', error);
            } finally {
                isLoading = false;
                document.getElementById('loadingIndicator').classList.remove('active');
            }
        }

        function createEntryCard(entry) {
            const card = document.createElement('div');
            card.className = 'entry-card';

            let titleHTML = '';
            if (entry.title) {
                titleHTML = '<div style="font-weight: 600; font-size: 16px; margin-bottom: 8px;">' + escapeHtml(entry.title) + '</div>';
            }

            let photoHTML = '';
            if (entry.photo) {
                photoHTML = '<img src="' + entry.photo + '" alt="Entry photo" class="entry-photo">';
            }

            card.innerHTML =
                '<div class="entry-header">' +
                    '<span class="entry-id">Entry #' + entry.id + '</span>' +
                    '<span class="entry-date">' + entry.createdAt + '</span>' +
                '</div>' +
                titleHTML +
                '<div class="entry-content">' + escapeHtml(entry.content) + '</div>' +
                photoHTML +
                '<div class="entry-actions">' +
                    '<button onclick="openEditModal(' + entry.id + ', \'' + escapeJs(entry.content) + '\')">Edit</button>' +
                    '<button class="btn-danger" onclick="openDeleteModal(' + entry.id + ')">Delete</button>' +
                '</div>';

            return card;
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        function escapeJs(text) {
            return text.replace(/\\/g, '\\\\').replace(/'/g, "\\'").replace(/"/g, '\\"').replace(/\n/g, '\\n').replace(/\r/g, '\\r');
        }

        window.addEventListener('scroll', () => {
            if (isLoading || !hasMore) return;
            const scrollPosition = window.innerHeight + window.scrollY;
            const pageHeight = document.documentElement.scrollHeight;
            if (pageHeight - scrollPosition < 200) {
                loadEntries();
            }
        });

        loadEntries();

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

                // Determine best supported format
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

                    // Stop all tracks in the stream
                    stream.getTracks().forEach(track => track.stop());
                };

                mediaRecorder.start();
                recordingStartTime = Date.now();

                // Update UI
                document.getElementById('startRecordBtn').style.display = 'none';
                document.getElementById('stopRecordBtn').style.display = 'inline-flex';
                document.getElementById('recordingIndicator').classList.add('active');
                document.getElementById('audioPreview').classList.remove('active');

                // Start timer
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

            // Update UI
            document.getElementById('startRecordBtn').style.display = 'inline-flex';
            document.getElementById('stopRecordBtn').style.display = 'none';

            // Stop timer
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

            // Reset UI
            document.getElementById('audioPreview').classList.remove('active');
            document.getElementById('audioPlayback').src = '';
            document.getElementById('recordingTime').textContent = '00:00';

            // Clear the file input
            const mediaInput = document.getElementById('media');
            mediaInput.value = '';
            document.getElementById('mediaFileName').textContent = 'Choose file';
        }

        function useRecording() {
            if (!recordedBlob) {
                alert('No recording available');
                return;
            }

            // Create a File object from the blob
            const fileName = 'recording-' + Date.now() + '.webm';
            const file = new File([recordedBlob], fileName, { type: recordedBlob.type });

            // Create a DataTransfer to set the file input
            const dataTransfer = new DataTransfer();
            dataTransfer.items.add(file);
            document.getElementById('media').files = dataTransfer.files;

            // Update file name display
            document.getElementById('mediaFileName').textContent = fileName;

            // Show confirmation
            document.getElementById('audioPreview').innerHTML = '<div style="color: #155724; font-weight: 600;">✓ Recording ready to upload</div>';
        }

        // Form submit handler - auto-stop recording and attach if in progress
        document.getElementById('createForm').addEventListener('submit', function(e) {
            const audioSourceRadio = document.querySelector('input[name="audio_source"]:checked');
            const isRecordMode = audioSourceRadio && audioSourceRadio.value === 'record';
            const isRecording = mediaRecorder && mediaRecorder.state === 'recording';
            const hasRecordedBlob = recordedBlob !== null;
            const fileInput = document.getElementById('media');
            const hasFileSelected = fileInput.files && fileInput.files.length > 0;

            // If in record mode and actively recording, stop and wait for blob
            if (isRecordMode && isRecording) {
                e.preventDefault();

                // Set up a one-time handler to submit after recording stops
                const originalOnStop = mediaRecorder.onstop;
                mediaRecorder.onstop = function(event) {
                    // Call original handler to create the blob
                    if (originalOnStop) {
                        originalOnStop.call(mediaRecorder, event);
                    }

                    // Wait a brief moment for blob to be ready, then attach and submit
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

                // Stop the recording
                mediaRecorder.stop();
                if (recordingTimer) {
                    clearInterval(recordingTimer);
                    recordingTimer = null;
                }
                return;
            }

            // If in record mode with a recorded blob but not yet attached, attach it
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
    </script>
</body>
</html>`
