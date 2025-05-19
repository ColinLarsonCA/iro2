import { AppShell } from "@mantine/core";
// import { useDisclosure } from "@mantine/hooks";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { CollabPage, HomePage } from "./pages";
import { TestPage } from "./pages/TestPage";

function App() {
  // const [opened, { toggle }] = useDisclosure();

  return (
    <AppShell
      // header={{ height: 60 }}
      // navbar={{
      //   width: 300,
      //   breakpoint: "sm",
      //   collapsed: { mobile: !opened },
      // }}
      padding="md"
    >
      {/* <AppShell.Header>
        <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
        <div>Logo</div>
      </AppShell.Header> */}
      {/* <AppShell.Navbar p="md">Navbar</AppShell.Navbar> */}
      <AppShell.Main>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/collab/:id" element={<CollabPage />} />
            <Route path="/test" element={<TestPage />} />
          </Routes>
        </BrowserRouter>
      </AppShell.Main>
    </AppShell>
  );
}

export default App;
