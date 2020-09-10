import React, { useState } from "react";
import "./App.css";
import "@rmwc/drawer/dist/styles";
import "@rmwc/list/dist/styles";
import "@rmwc/button/dist/styles";
import "@rmwc/icon/dist/styles";
import "@rmwc/elevation/dist/styles";
import { Drawer, DrawerContent, DrawerHeader } from "@rmwc/drawer";
import { List, ListItem } from "@rmwc/list";
import { Icon } from "@rmwc/icon";
import { Elevation } from "@rmwc/elevation";
import Logo from "./logo-white.svg";

function App() {
  return (
    <div className="App">
      <Elevation z={24}>
        <Drawer
          style={{
            backgroundColor: "rgba(57, 156, 248, 1)",
          }}
        >
          <DrawerHeader>
            <img src={Logo} alt={"Timeterm Logo"} width={96} />
          </DrawerHeader>
          <DrawerContent>
            <List>
              <ListItem style={{ color: "white" }}>
                <Icon icon="tablet" />
                &nbsp; Apparaten
              </ListItem>
              <ListItem style={{ color: "white" }}>
                <Icon icon="group" />
                &nbsp; Gebruikers
              </ListItem>
              <ListItem style={{ color: "white" }}>
                <Icon icon="bluetooth" />
                &nbsp; Apparaat koppelen
              </ListItem>
            </List>
          </DrawerContent>
        </Drawer>
      </Elevation>
    </div>
  );
}

export default App;
