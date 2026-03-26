import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import App from "./App.tsx";
import { CookiesProvider } from "react-cookie";

const router = createBrowserRouter([
  { path: "/", element: <App roomSlug="/" /> },
  { path: "/residence-hall-a", element: <App roomSlug="residence-hall-a" /> },
  { path: "/residence-hall-b", element: <App roomSlug="residence-hall-b" /> },
  { path: "/residence-hall-c", element: <App roomSlug="residence-hall-c" /> },
  { path: "/kate-gleason", element: <App roomSlug="kate-gleason" /> },
  { path: "/gibson", element: <App roomSlug="gibson" /> },
  { path: "/peterson", element: <App roomSlug="peterson" /> },
  { path: "/sol-heumann", element: <App roomSlug="sol-heumann" /> },
]);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <CookiesProvider defaultSetOptions={{ path: "/" }}>
      <RouterProvider router={router} />
    </CookiesProvider>
  </StrictMode>,
);
