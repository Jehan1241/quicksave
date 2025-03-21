export const handleApiError = async (response: Response) => {
  let errorMessage = `HTTP Error: ${response.status} - ${response.statusText}`;
  let errorDetails = "";

  if (response.status === 404) {
    errorMessage = "Route not found (404)";
  } else {
    try {
      const errorResp = await response.json();
      errorMessage = errorResp.error || errorMessage;
      errorDetails = errorResp.details || "";
    } catch {
      errorMessage = "Failed to parse server error response";
    }
  }
  throw new Error(`${errorMessage} ${errorDetails}`.trim());
};
