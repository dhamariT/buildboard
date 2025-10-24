/**
 * API client for BuildBoard backend
 */

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

/**
 * Base fetch wrapper with error handling
 */
async function apiFetch(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;

  const config = {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  };

  try {
    const response = await fetch(url, config);
    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || `HTTP error! status: ${response.status}`);
    }

    return data;
  } catch (error) {
    console.error('API Error:', error);
    throw error;
  }
}

/**
 * Early Start API endpoints
 */
export const earlyStartAPI = {
  /**
   * Sign up for early access
   * @param {Object} data - Signup data
   * @param {string} data.email - User email
   * @param {string} [data.firstName] - User first name
   * @param {string} [data.lastName] - User last name
   * @returns {Promise<Object>} Response with message
   */
  signup: async (data) => {
    const response = await apiFetch('/early-start/signup', {
      method: 'POST',
      body: JSON.stringify(data),
    });

    // In development mode, log the OTP to console
    const isDevelopment = process.env.NEXT_PUBLIC_ENV === 'development';
    if (isDevelopment && response.otp) {
      console.log('🔐 Development Mode - OTP Code:', response.otp);
    }

    return response;
  },

  /**
   * Verify OTP code
   * @param {Object} data - Verification data
   * @param {string} data.email - User email
   * @param {string} data.otp - 6-character OTP code
   * @returns {Promise<Object>} Response with verification status
   */
  verifyOTP: async (data) => {
    return apiFetch('/early-start/verify', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  /**
   * Get signup count statistics
   * @returns {Promise<Object>} Object with total and verified counts
   */
  getCount: async () => {
    return apiFetch('/early-start/count', {
      method: 'GET',
    });
  },
};

/**
 * Admin API endpoints
 */
export const adminAPI = {
  /**
   * Get list of all early start signups (admin only)
   * @returns {Promise<Object>} Object with users array and pagination info
   */
  listSignups: async () => {
    return apiFetch('/admin/early-start', {
      method: 'GET',
    });
  },
};

/**
 * Health check endpoint
 */
export const healthCheck = async () => {
  return apiFetch('/health', {
    method: 'GET',
  });
};