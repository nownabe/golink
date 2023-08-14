import { createBrowserRouter } from "react-router-dom";

import Layout, { layoutLoader } from "./pages/Layout";
import Home from "./pages/Home";
import MyGolinks, { myGolinksLoader } from "./pages/MyGolinks";
import NewGolink, { newGolinkLoader } from "./pages/NewGolink";
import EditGolink, { editGolinkLoader } from "./pages/EditGolink";
import RouterErrorDialog from "./components/RouterErrorDialog";

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
          errorElement: <RouterErrorDialog />,
        },
        {
          path: "/-/new",
          element: <NewGolink />,
          loader: newGolinkLoader,
          errorElement: <RouterErrorDialog />,
        },
        {
          path: "/:name",
          element: <EditGolink />,
          loader: editGolinkLoader,
          errorElement: <RouterErrorDialog />,
        },
      ],
    },
  ],
  {
    basename: "/c/",
  }
);

export default router;
