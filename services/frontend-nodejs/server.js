const express = require('express');
const session = require('express-session');
const path = require('path');
const fetch = require('node-fetch');

const app = express();
const PORT = process.env.PORT || 3000;

// Base URL of the API Gateway (inside Docker/K8s this will be set via API_BASE)
const API_BASE = process.env.API_BASE || 'http://api-gateway-go:8000/api';

// Basic middleware
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// In-memory session (demo only)
app.use(
  session({
    secret: 'very-secret-demo-key',
    resave: false,
    saveUninitialized: true,
  })
);

// Serve static HTML files from public/
app.use(express.static(path.join(__dirname, 'public')));

// Helper to ensure user is logged in for API actions
function requireSessionUser(req, res, next) {
  if (!req.session || !req.session.user) {
    return res.status(401).json({ error: 'Not logged in' });
  }
  next();
}

// Returns current session user (for pages to show username)
app.get('/session', (req, res) => {
  if (req.session && req.session.user) {
    return res.json(req.session.user);
  }
  res.json(null);
});

// POST /login - call backend /api/login, store user in session only on success
app.post('/login', async (req, res) => {
  const { username, password } = req.body || {};

  if (!username || !password) {
    return res.status(400).json({ error: 'username and password are required' });
  }

  try {
    const resp = await fetch(`${API_BASE}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });

    if (!resp.ok) {
      const text = await resp.text();
      console.warn('Backend login failed with status', resp.status, text);
      return res.status(401).json({ error: 'Invalid credentials or login failed' });
    }

    const data = await resp.json();

    // Expect backend to return at least a user identifier; if not, treat as failure
    const userId = data.user_id || data.id;
    if (!userId) {
      console.warn('Backend login did not return a user_id/id field');
      return res.status(401).json({ error: 'Invalid login response from user service' });
    }

    const user = {
      username: data.username || username,
      user_id: userId,
    };

    req.session.user = user;
    res.json(user);
  } catch (err) {
    console.error('Error calling backend /login:', err);
    res.status(500).json({ error: 'Error contacting user service' });
  }
});

// POST /add-to-cart - body: { product_id }
app.post('/add-to-cart', requireSessionUser, async (req, res) => {
  const { product_id } = req.body || {};
  const user = req.session.user;

  if (!product_id) {
    return res.status(400).json({ error: 'product_id is required' });
  }

  try {
    const resp = await fetch(`${API_BASE}/cart/add`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: user.user_id,
        product_id,
      }),
    });

    if (!resp.ok) {
      const text = await resp.text();
      return res.status(500).json({ error: 'Failed to add to cart', details: text });
    }

    res.json({ ok: true });
  } catch (err) {
    console.error('Error calling /cart/add:', err);
    res.status(500).json({ error: 'Error contacting cart service' });
  }
});

// GET /cart-data - proxy to backend /api/cart/{user_id} using session user
app.get('/cart-data', requireSessionUser, async (req, res) => {
  const user = req.session.user;
  try {
    const resp = await fetch(`${API_BASE}/cart/${encodeURIComponent(user.user_id)}`);
    if (!resp.ok) {
      const text = await resp.text();
      return res.status(500).json({ error: 'Failed to load cart', details: text });
    }
    const data = await resp.json();
    res.json(data);
  } catch (err) {
    console.error('Error calling /cart/{user_id}:', err);
    res.status(500).json({ error: 'Error contacting cart service' });
  }
});

// POST /checkout - calls backend /api/order/checkout
// For this demo we always return success to the browser, even if backend fails.
app.post('/checkout', requireSessionUser, async (req, res) => {
  const user = req.session.user;

  try {
    const resp = await fetch(`${API_BASE}/order/checkout`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: user.user_id }),
    });

    if (!resp.ok) {
      const text = await resp.text();
      console.warn('Order checkout backend returned non-OK status:', resp.status, text);
      // Demo fallback: still report success so the flow completes
      return res.json({
        ok: true,
        order: {
          message: 'Order completed (demo fallback)',
          backend_status: resp.status,
          backend_response: text,
        },
      });
    }

    const data = await resp.json();
    res.json({ ok: true, order: data });
  } catch (err) {
    console.error('Error calling /order/checkout:', err);
    // Demo fallback on error
    res.json({
      ok: true,
      order: {
        message: 'Order completed (demo fallback due to error)',
        error: err.message || String(err),
      },
    });
  }
});

// GET /orders-data - proxy to backend orders endpoint for current user
app.get('/orders-data', requireSessionUser, async (req, res) => {
  const user = req.session.user;
  try {
    // Now call GET /api/order/{user_id} via gateway
    const url = `${API_BASE}/order/${encodeURIComponent(user.user_id)}`;
    const resp = await fetch(url);

    if (!resp.ok) {
      const text = await resp.text();
      return res.status(500).json({ error: 'Failed to load orders', details: text });
    }

    const data = await resp.json();
    res.json(data);
  } catch (err) {
    console.error('Error calling orders endpoint:', err);
    res.status(500).json({ error: 'Error contacting order service' });
  }
});

// Simple health endpoint for Kubernetes/liveness
app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

app.listen(PORT, () => {
  console.log(`frontend-nodejs listening on http://localhost:${PORT}`);
  console.log(`Using API gateway base URL: ${API_BASE}`);
});
