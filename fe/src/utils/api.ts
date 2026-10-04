// Utility to resolve the base API endpoint URL
// In browsers, relative '/api/v1' hits the current host and Next.js internal proxy
export const getApiUrl = (): string => {
  return process.env.NEXT_PUBLIC_API_URL || '/api/v1';
};
