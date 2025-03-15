self.onmessage = async (e) => {
  const { url, cacheBuster } = e.data;

  try {
    const response = await fetch(url, { method: "HEAD" });

    if (response.ok) {
      self.postMessage({ url: url + "?t=" + cacheBuster, success: true });
    } else {
      self.postMessage({ url, success: false });
    }
  } catch {
    self.postMessage({ url, success: false });
  }
};
