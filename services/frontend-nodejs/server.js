const express = require('express');
const http = require('http');

const app = express();
const PORT = process.env.PORT || 3000;

// Frontend will call /api/* on the same origin (localhost:3000)
// This Express app will proxy /api/* to the Go API gateway on localhost:8000
const API_BASE = '/api';
const GATEWAY_BASE = 'http://localhost:8000';

// Basic Express setup
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Very simple proxy: /api/* -> http://localhost:8000/api/*
app.use('/api', (req, res) => {
  const targetUrl = GATEWAY_BASE + req.originalUrl; // e.g. /api/products
  const options = new URL(targetUrl);

  options.method = req.method;
  options.headers = req.headers;

  const proxyReq = http.request(options, (proxyRes) => {
    res.writeHead(proxyRes.statusCode || 500, proxyRes.headers);
    proxyRes.pipe(res, { end: true });
  });

  proxyReq.on('error', (err) => {
    console.error('Proxy error:', err);
    res.status(502).json({ error: 'Bad gateway' });
  });

  if (req.readable) {
    req.pipe(proxyReq, { end: true });
  } else {
    proxyReq.end();
  }
});

// Small helper to render a full HTML page
function renderPage(title, bodyHtml) {
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>${title}</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 0; padding: 0; }
    header { background: #222; color: #fff; padding: 10px 16px; }
    main { padding: 16px; }
    .container { max-width: 800px; margin: 0 auto; }
    .error { color: #c00; margin-top: 8px; }
    .success { color: #060; margin-top: 8px; }
    .product-list { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 12px; }
    .product { border: 1px solid #ddd; padding: 8px; border-radius: 4px; }
    .product h3 { margin: 0 0 4px; font-size: 16px; }
    nav a { color: #fff; margin-right: 12px; text-decoration: none; }
    nav a:hover { text-decoration: underline; }
    button { padding: 4px 8px; cursor: pointer; }
    table { border-collapse: collapse; width: 100%; margin-top: 8px; }
    th, td { border: 1px solid #ddd; padding: 6px 8px; text-align: left; }
    th { background: #f4f4f4; }
    form { margin-top: 8px; }
    input[type="text"], input[type="password"] { padding: 4px 6px; width: 200px; }
  </style>
</head>
<body>
  <header>
    <div class="container">
      <nav>
        <a href="/login">Login</a>
        <a href="/home">Home</a>
        <a href="/cart">Cart</a>
      </nav>
    </div>
  </header>
  <main>
    <div class="container">
      ${bodyHtml}
    </div>
  </main>
</body>
</html>`;
}

// Redirect root to /login for simplicity
app.get('/', (req, res) => {
  res.redirect('/login');
});

// /login – accept any credentials, but still try backend if available
app.get('/login', (req, res) => {
  const body = `
    <h1>Login</h1>
    <p>This demo accepts any username/password. It will also call the backend login API if available.</p>
    <form id="login-form">
      <div>
        <label>Username:</label><br />
        <input type="text" id="username" required />
      </div>
      <div style="margin-top:8px;">
        <label>Password:</label><br />
        <input type="password" id="password" required />
      </div>
      <div style="margin-top:12px;">
        <button type="submit">Login</button>
      </div>
      <div id="message" class="error"></div>
    </form>
    <script>
      const API_BASE = '${API_BASE}';
      document.getElementById('login-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        // Accept any credentials by default
        let user = { username, user_id: 'demo-user' };

        try {
          const resp = await fetch(API_BASE + '/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
          });
          if (resp.ok) {
            const data = await resp.json();
            user = data;
          }
        } catch (err) {
          console.warn('Backend login failed, using demo user', err);
        }

        localStorage.setItem('user', JSON.stringify(user));
        window.location.href = '/home';
      });
    </script>
  `;
  res.send(renderPage('Login', body));
});

// /home – list products from backend
app.get('/home', (req, res) => {
  const body = `
    <h1>Products</h1>
    <div id="products" class="product-list"></div>
    <div id="message" class="error"></div>
    <script>
      const API_BASE = '${API_BASE}';

      async function loadProducts() {
        try {
          const url = API_BASE + '/products';
          console.log('Fetching products from:', url);
          const resp = await fetch(url);
          console.log('Products response status:', resp.status);

          const text = await resp.text();
          console.log('Raw products response:', text);

          if (!resp.ok) {
            document.getElementById('message').textContent = 'Failed to load products';
            return;
          }

          let products = [];
          try {
            products = JSON.parse(text);
          } catch (parseErr) {
            console.error('Failed to parse products JSON:', parseErr);
            document.getElementById('message').textContent = 'Invalid products data from server';
            return;
          }

          const container = document.getElementById('products');
          container.innerHTML = '';

          if (!Array.isArray(products) || !products.length) {
            container.textContent = 'No products available';
            return;
          }

          products.forEach((p) => {
            const div = document.createElement('div');
            div.className = 'product';

            let html = '';
            const name = p.name || p.id;
            html += '<h3>' + name + '</h3>';
            html += '<div>Price: $' + p.price + '</div>';
            html += '<div style="margin-top:4px;">';
            html += '<label>Quantity:</label> ';
            html += '<input type="number" min="1" value="1" data-id="' + p.id + '" style="width:60px;" /> ';
            html += '<button data-id="' + p.id + '">Add to cart</button>';
            html += '</div>';

            div.innerHTML = html;

            const qtyInput = div.querySelector('input');
            const button = div.querySelector('button');
            button.addEventListener('click', () => {
              const qty = parseInt(qtyInput.value, 10) || 1;
              addToCart(p.id, qty);
            });

            container.appendChild(div);
          });
        } catch (err) {
          console.error('Error loading products:', err);
          document.getElementById('message').textContent = 'Error contacting server';
        }
      }

      async function addToCart(productId, quantity) {
        const user = JSON.parse(localStorage.getItem('user') || '{}');
        const userId = user.user_id || 'demo-user';
        try {
          for (let i = 0; i < quantity; i++) {
            const resp = await fetch(API_BASE + '/cart/add', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ user_id: userId, product_id: productId })
            });
            if (!resp.ok) {
              document.getElementById('message').textContent = 'Failed to add to cart';
              return;
            }
          }
          window.location.href = '/cart';
        } catch (err) {
          console.error('Error adding to cart:', err);
          document.getElementById('message').textContent = 'Error contacting server';
        }
      }

      loadProducts();
    </script>
  `;
  res.send(renderPage('Home', body));
});

// /cart – show items in cart
app.get('/cart', (req, res) => {
  const body = `
    <h1>Your Cart</h1>
    <div id="cart-container"></div>
    <div id="message" class="error"></div>
    <div style="margin-top:12px;">
      <a href="/checkout"><button>Proceed to Checkout</button></a>
    </div>
    <script>
      const API_BASE = '${API_BASE}';
      async function loadCart() {
        const user = JSON.parse(localStorage.getItem('user') || '{}');
        const userId = user.user_id || 'demo-user';
        try {
          const resp = await fetch(API_BASE + '/cart/' + encodeURIComponent(userId));
          if (!resp.ok) {
            document.getElementById('message').textContent = 'Failed to load cart';
            return;
          }
          const data = await resp.json();
          const items = data.items || [];
          const container = document.getElementById('cart-container');
          if (!items.length) {
            container.textContent = 'Cart is empty';
            return;
          }
          let html = '<table><thead><tr><th>Product ID</th><th>Qty</th></tr></thead><tbody>';
          items.forEach((it) => {
            html += '<tr><td>' + it.product_id + '</td><td>' + it.quantity + '</td></tr>';
          });
          html += '</tbody></table>';
          container.innerHTML = html;
        } catch (err) {
          console.error('Error loading cart:', err);
          document.getElementById('message').textContent = 'Error contacting server';
        }
      }
      loadCart();
    </script>
  `;
  res.send(renderPage('Cart', body));
});

// /checkout – trigger order checkout
app.get('/checkout', (req, res) => {
  const body = `
    <h1>Checkout</h1>
    <p>Click the button below to place your order.</p>
    <button id="checkout-btn">Checkout</button>
    <div id="message" class="error"></div>
    <div id="success" class="success"></div>
    <script>
      const API_BASE = '${API_BASE}';
      document.getElementById('checkout-btn').addEventListener('click', async () => {
        const user = JSON.parse(localStorage.getItem('user') || '{}');
        const userId = user.user_id || 'demo-user';
        try {
          const resp = await fetch(API_BASE + '/order', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ user_id: userId })
          });
          if (!resp.ok) {
            document.getElementById('message').textContent = 'Checkout failed';
            return;
          }
          const data = await resp.json();
          document.getElementById('success').textContent = 'Order successful! Payment ID: ' + (data.payment_id || 'n/a');
          setTimeout(() => { window.location.href = '/success'; }, 800);
        } catch (err) {
          console.error('Error during checkout:', err);
          document.getElementById('message').textContent = 'Error contacting server';
        }
      });
    </script>
  `;
  res.send(renderPage('Checkout', body));
});

// /success – simple success page
app.get('/success', (req, res) => {
  const body = `
    <h1>Success</h1>
    <p>Your order was placed successfully.</p>
    <p><a href="/home">Back to Home</a></p>
  `;
  res.send(renderPage('Success', body));
});

app.listen(PORT, () => {
  console.log(`Frontend listening on http://localhost:${PORT}`);
});
