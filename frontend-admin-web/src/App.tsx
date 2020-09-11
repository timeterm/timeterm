import React, { useState } from "react";
import "./App.css";
import { Drawer, DrawerContent, DrawerHeader } from "@rmwc/drawer";
import {
  List,
  ListItem,
  ListItemGraphic,
  ListItemPrimaryText,
  ListItemSecondaryText,
  ListItemText,
} from "@rmwc/list";
import { Elevation } from "@rmwc/elevation";
import { ThemeProvider, Theme } from "@rmwc/theme";
import Logo from "./logo-white.svg";
import DevicesTable, { DeviceStatus } from "./DevicesTable";
import { Button } from "@rmwc/button";

function App() {
  const [selectedItems, setSelectedItems] = useState([] as string[]);

  return (
    <ThemeProvider
      options={{
        primary: "rgba(57, 156, 248, 1)",
        secondary: "rgb(127,193,255)",
        onPrimary: "white",
      }}
      style={{
        height: "100%",
      }}
    >
      <div className="App">
        <Elevation
          z={24}
          style={{
            height: "100%",
          }}
        >
          <Theme use={["primaryBg", "onPrimary"]} wrap>
            <Drawer>
              <DrawerHeader>
                <img src={Logo} alt={"Timeterm Logo"} width={96} />
              </DrawerHeader>
              <DrawerContent
                style={{
                  display: "flex",
                  flexDirection: "column",
                }}
              >
                <List style={{ flexGrow: 1 }}>
                  <Theme use={["primaryBg", "onPrimary"]} wrap>
                    <ListItem>
                      <Theme use={["primaryBg", "onPrimary"]} wrap>
                        <ListItemGraphic icon="tablet" />
                      </Theme>
                      Apparaten
                    </ListItem>
                  </Theme>
                  <Theme use={["primaryBg", "onPrimary"]} wrap>
                    <ListItem>
                      <Theme use={["primaryBg", "onPrimary"]} wrap>
                        <ListItemGraphic icon="group" />
                      </Theme>
                      Gebruikers
                    </ListItem>
                  </Theme>
                  <Theme use={["primaryBg", "onPrimary"]} wrap>
                    <ListItem>
                      <Theme use={["primaryBg", "onPrimary"]} wrap>
                        <ListItemGraphic icon="bluetooth_connected" />
                      </Theme>
                      Apparaat koppelen
                    </ListItem>
                  </Theme>
                </List>

                <List twoLine={true} style={{ paddingBottom: 0 }}>
                  <Theme use={["primaryBg", "onPrimary"]} wrap>
                    <ListItem>
                      <Theme use={["primaryBg", "onPrimary"]} wrap>
                        <ListItemGraphic icon="person" />
                      </Theme>
                      <ListItemText>
                        <ListItemPrimaryText>Admin</ListItemPrimaryText>
                        <Theme use={["primaryBg", "onPrimary"]} wrap>
                          <ListItemSecondaryText>
                            admin@timeterm.nl
                          </ListItemSecondaryText>
                        </Theme>
                      </ListItemText>
                    </ListItem>
                  </Theme>
                </List>
                <List>
                  <Theme use={["primaryBg", "onPrimary"]} wrap>
                    <ListItem>
                      <Theme use={["primaryBg", "onPrimary"]} wrap>
                        <ListItemGraphic icon="logout" />
                      </Theme>
                      Uitloggen
                    </ListItem>
                  </Theme>
                </List>
              </DrawerContent>
            </Drawer>
          </Theme>
        </Elevation>

        <div
          style={{
            display: "flex",
            flexDirection: "column",
            width: "100%",
          }}
        >
          <div
            style={{
              display: "flex",
              marginLeft: 32,
              marginTop: 16,
              marginRight: 16,
              height: 40,
              justifyContent: "space-between",
            }}
          >
            <h1 style={{ marginTop: 0 }}>Apparaten</h1>
            <div>
              <Button
                icon={"delete"}
                danger
                raised
                disabled={selectedItems.length === 0}
              >
                Verwijderen
              </Button>
              <Button
                icon={"power_settings_new"}
                danger
                disabled={selectedItems.length === 0}
                raised
                style={{ marginLeft: 8 }}
              >
                Opnieuw opstarten
              </Button>
            </div>
          </div>

          <Elevation
            z={16}
            style={{
              flexGrow: 1,
              margin: 16,
              borderRadius: 8,
            }}
          >
            <DevicesTable
              devices={[
                {
                  name: "Mediatheek 1",
                  id: "192803410234",
                  status: DeviceStatus.Online,
                },
                {
                  name: "Mediatheek 2",
                  id: "928034810234",
                  status: DeviceStatus.Offline,
                },
              ]}
              setSelectedItems={setSelectedItems}
            />
          </Elevation>
        </div>
      </div>
    </ThemeProvider>
  );
}

export default App;
