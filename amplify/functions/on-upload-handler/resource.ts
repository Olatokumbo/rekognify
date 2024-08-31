import { defineFunction } from '@aws-amplify/backend';
import { projectName } from '../../constant';

export const onUploadHandler = defineFunction({
    name: `${projectName}`,
    entry: "./index.ts",
})