import { defineStorage } from '@aws-amplify/backend';
import { projectName } from '../constant';

export const storage = defineStorage({
    name: `${projectName}-storage`,
});