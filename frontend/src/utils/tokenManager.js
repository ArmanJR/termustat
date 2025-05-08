let accessToken = null;
let setAccessTokenCallback = null;
let tryRefreshTokenCallback = null;

export const tokenManager = {
  get: () => accessToken,
  set: (token) => {
    accessToken = token;
    if (setAccessTokenCallback) setAccessTokenCallback(token);
  },
  onSet: (callback) => {
    setAccessTokenCallback = callback;
  },
  onRefresh: (callback) => {
    tryRefreshTokenCallback = callback;
  },
  tryRefresh: async () => {
    if (tryRefreshTokenCallback) {
      return await tryRefreshTokenCallback();
    }
    return null;
  },
  clear: () => {
    accessToken = null;
  },
};
