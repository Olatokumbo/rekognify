import { defineStorage } from '@aws-amplify/backend';
import { projectName } from '../constant';
import { onUploadHandler } from '../functions/on-upload-handler/resource';

export const storage = defineStorage({
    name: `${projectName}-storage`,
    triggers: {
        onUpload: onUploadHandler
    }
});