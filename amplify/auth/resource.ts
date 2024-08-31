import { defineAuth } from '@aws-amplify/backend';
import { projectName } from '../constant';

/**
 * Define and configure your auth resource
 * @see https://docs.amplify.aws/gen2/build-a-backend/auth
 */
export const auth = defineAuth({
  name: `${projectName}-auth`,
  loginWith: {
    email: true,
  },
});
