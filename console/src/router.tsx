import { createBrowserRouter } from "react-router-dom";

import Layout, { layoutLoader } from "./pages/Layout";
import Home from "./pages/Home";
import MyGolinks, { myGolinksLoader } from "./pages/MyGolinks";
import NewGolink, { newGolinkLoader } from "./pages/NewGolink";
import EditGolink, { editGolinkLoader } from "./pages/EditGolink";
import ErrorDialog from "./components/ErrorDialog";

const router = createBrowserRouter(
  [
    {
      path: "/",
      element: <Layout />,
      loader: layoutLoader,
      errorElement: <ErrorDialog />,
      children: [
        {
          index: true,
          element: <Home />,
        },
        {
          path: "/-/",
          element: <MyGolinks />,
          loader: myGolinksLoader,
          errorElement: <ErrorDialog />,
        },
        {
          path: "/-/new",
          element: <NewGolink />,
          loader: newGolinkLoader,
          errorElement: <ErrorDialog />,
        },
        {
          path: "/:name",
          element: <EditGolink />,
          loader: editGolinkLoader,
          errorElement: <ErrorDialog />,
        },
      ],
    },
  ],
  {
    basename: "/c/",
  }
);

export default router;
