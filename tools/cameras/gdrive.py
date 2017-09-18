import os
import logging

# Public class for the most common GDrive operations needed (currently just pushFile())
class GDrive:
    def __init__(self, gDriveCmd = "/usr/local/bin/gdrive", cacheLocation = os.path.join(os.path.expanduser("~"), ".gdrive_cache")):
        self._impl = _Impl(gDriveCmd)
        self._cache = _IdCache(cacheLocation, self._impl)  

    def pushFile(self, filename, remotePath):
        logging.debug("++ Uploading file: " + filename + " to " + remotePath)
        parentId = self._cache.getId(remotePath)
        self._impl.uploadFile(filename, parentId)


# The GDrive implementation class - deals with GDrive only
class _Impl:
    def __init__(self, gDriveCmd):
        self.gDriveCmd = gDriveCmd

    def createFolder(self, folderName, parentId):
        # returns the id of the new folder
        parentOpt = "" if parentId == "" else "-p {id}".format(id = parentId)
        fullCmd = "{cmd} mkdir {opt} \"{name}\" | sed s/Directory// | sed s/created//".format(cmd = self.gDriveCmd, opt = parentOpt, id = parentId, name = folderName)
        id = os.popen(fullCmd).read().strip()
        logging.debug("  -> createFolder: '{}' (parent id: '{}' -> id: '{}'".format(folderName, parentId, id))
        return id

    def uploadFile(self, filename, parentId):
        parentOpt = "" if parentId == "" else "-p {id}".format(id = parentId)
        fullCmd = "{cmd} upload {opt} \"{file}\"".format(cmd = self.gDriveCmd, opt = parentOpt, file = filename)
        os.popen(fullCmd).read()
        logging.debug("  -> uploadFile: '{}' (parent id: '{}'".format(filename, parentId))

    def getId(self, name, parentId):
        # this will return the ID of the most recent folder with this name
        parentOpt = "" if parentId == "" else "-q \"'{id}' in parents\"".format(id = parentId)
        fullCmd = "{cmd} list {opt} --order folder | grep \"{name}\" | head -1 | sed s/\\\\s.*//".format(cmd = self.gDriveCmd, opt = parentOpt, id = parentId, name = name)
        id = os.popen(fullCmd).read().strip()
        if(id == ""): 
            raise IOError("File or folder with name \"{name}\" and parent id \"{id}\" not found on Google Drive".format(name = name, id=parentId))
        logging.debug("  -> getId: '{}' (parent id: '{}' -> id: '{}'".format(name, parentId, id))
        return id


# The class that caches folder IDs locally to speed up lookups  
class _IdCache:
    def __init__(self, location, gdrive):
        self.location = location
        self.gdrive = gdrive

        if not os.path.exists(location):
            os.makedirs(location)
    
    def _getCacheFilenameFromPath(self, path):
        return os.path.join(self.location, path) + ".id" 

    # retrieves the ID of a given directory specified by path - if the dir does not exist
    # on Google drive, it will be created first
    def getId(self, path):
        if path == "":
            raise ValueError("Empty path.")

        # try the local cache first
        cacheFilename = self._getCacheFilenameFromPath(path)
        if os.path.exists(cacheFilename):
            return open(cacheFilename, 'r').read()

        
        # split path into <parent> and <dir/file>
        #filename = fullPath[fullPath.rfind('/')+1:]
        #parentPath = fullPath[:fullPath.rfind('/')-1]
        (parentPath, dirname) = os.path.split(path)

        logging.debug("parent: {}, dirname: {}".format(parentPath, dirname))

        # create hierarchy to avoid creating each path individually later
        if parentPath != "" and not os.path.exists(os.path.join(self.location, parentPath)):
            os.makedirs(os.path.join(self.location, parentPath))

        # we need the id of the parent path (if we have a parent)
        parentId = ""
        if parentPath != "":
            # parentCacheFilename = self._getCacheFilenameFromPath(parentPath)

            # if os.path.exists(parentCacheFilename):
            #     parentId = open(parentCacheFilename, 'r').read().strip() 
            # else:    
                parentId = self.getId(parentPath)
        
        # we don't have the id of the requested file/dir - get it from gdrive
        id = ""
        try:
            id = self.gdrive.getId(dirname, parentId)
        except IOError:
            # we need to create it first
            id = self.gdrive.createFolder(dirname, parentId)

        # cache the result
        cacheFile = open(cacheFilename, "w")
        cacheFile.write(id)
        cacheFile.close()

        return id
