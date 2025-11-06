export const handleApiError = (e: any, thunkAPI: any) => {
  if (e && e.response) {
    const errorMessage = e?.response?.data?.message;
    const errorName = e?.response?.data?.error;
    if (Array.isArray(errorMessage)) {
      const payload = {
        message: errorName,
        description: errorMessage
          .map((item, index) => `${index + 1}) ${item}`)
          .join(" "),
      };
      return thunkAPI.rejectWithValue(payload);
    } else if (errorMessage || errorName) {
      const payload = {
        message: errorName,
        description: errorMessage,
      };
      return thunkAPI.rejectWithValue(payload);
    } else {
      const payload = {
        message: errorName,
        description: e?.response?.statusText,
      };

      return thunkAPI.rejectWithValue(payload);
    }
  }
};
