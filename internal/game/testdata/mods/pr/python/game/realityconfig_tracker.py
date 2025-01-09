# ========================================================================================================
#
# PROJECT REALITY SERVER SETTINGS DEFINITION FOR THE REALITYTRACKER SCRIPT
#
# This file can be fully edited and it's automatically used by local and single-player servers.
#
#
#
C = {}

# Set to false to completely disable the tracker
C['ENABLE'] = True

# TRACKER UPDATE INTERVAL
# Every [UPDATE_INTERVAL] the server calls an update that function that collects all the relevant
# data from the server and writes it to a file
C['UPDATE_INTERVAL'] = 0.3

# ================= Local work mode settings

# Folder to write incomplete recordings into. Keep folder private to prevent ghosting!
C['TMP_FOLDER'] = 'temp'

# FILE NAME
# available parameters:
# - '/map' '/mode' '/layer'
# - date related strings that are parsed by strftime (
# https://docs.python.org/2/library/datetime.html#strftime-strptime-behavior)
C['FILE_NAME'] = 'tracker_%Y_%m_%d_%H_%M_%S_/map_/mode_/layer'

# ========PUBLIC TRACKER FILE
C['TRACKER_FILE_PUBLIC'] = True

# Folder to move complete public recordings into.
C['PUBLIC_FOLDER'] = 'demos'

# Public file private data selection:

# Write player's IP to the public file.
C['FILE_PRIVATEDATA_IP'] = False
# Write player's HASH to the public file.
C['FILE_PRIVATEDATA_HASH'] = True

# Enable any chat recording
C['CHAT_ENABLE'] = True
# Enable Team chat recording
C['CHAT_TEAM'] = True
# Enable squad chat recording
C['CHAT_SQUAD'] = False

# ==== PRIVATE TRACKER FILE
# Create an extra file without filtering any private information
C['TRACKER_FILE_PRIVATE'] = False
C['PRIVATE_FOLDER'] = 'demos_private'

# ===== JSON Summary
# Write a summary at end of round.
C['JSON_ENABLE'] = True
# Folder for the JSON files. This must be set to something if JSON is enabled.
C['JSON_FOLDER'] = 'json'

C['JSON_WRITE_IP'] = True
C['JSON_WRITE_HASH'] = True

# ===== Advanced options
# Flush to file every recording tick, useful if you're having another program read the file.
C['OUTPUT_FLUSH_EVERY_TICK'] = False

# Work in progress, doesn't work:

# ===== Networking work mode settings
# TRACKER NETWORKING
# Specify if you want to run the tracker with a remote connection
C['TRACKER_NETWORKING'] = False

# TRACKER TCP SERVER PORT
# Edit this setting to set the port that the server will listen to.
# If you change this setting to anything other than None it will be used instead.
C['SERVER_PORT'] = 6669
