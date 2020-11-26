import { Typography } from "@rmwc/typography";
import { TextField } from "@rmwc/textfield";
import { Select } from "@rmwc/select";
import { Switch } from "@rmwc/switch";
import React from "react";
import { CollapsibleList, List, SimpleListItem } from "@rmwc/list";
import { fetchAuthnd } from "../DevicesPage";
import { Paginated } from "../GeneralTable";
import { Button } from "@rmwc/button";
import useSetting, { SettingPageProps } from "./useSetting";
import { Icon } from "@rmwc/icon";

enum NetworkingServiceType {
  Ethernet = "Ethernet",
  Wifi = "Wifi",
}

enum Ipv4ConfigType {
  Off = "Off",
  Dhcp = "Dhcp",
  Custom = "Custom",
}

enum Ipv6ConfigType {
  Off = "Off",
  Auto = "Auto",
  Custom = "Custom",
}

enum Ipv6Privacy {
  Disabled = "Disabled",
  Enabled = "Enabled",
  Preferred = "Preferred",
}

enum Security {
  Psk = "Psk",
  Ieee8021x = "ieee8021x",
  None = "None",
  Wep = "Wep",
}

enum EapType {
  Tls = "Tls",
  Ttls = "Ttls",
  Peap = "Peap",
}

enum CaCertType {
  Pem = "Pem",
  Der = "Der",
}

enum PrivateKeyType {
  Pem = "Pem",
  Der = "Der",
  Pfx = "Pfx",
}

enum PrivateKeyPassphraseType {
  Fsid = "Fsid",
}

enum Phase2Type {
  Gtc = "Gtc",
  MschapV2 = "MschapV2",
}

interface Ipv4ConfigSettings {
  network?: string;
  netmask?: string;
  gateway?: string;
}

interface Ipv4Config {
  type?: Ipv4ConfigType;
  settings?: Ipv4ConfigSettings;
}

interface Ipv6ConfigSettings {
  network?: string;
  prefixLength?: number;
  gateway?: string;
}

interface Ipv6Config {
  type?: Ipv6ConfigType;
  settings?: Ipv6ConfigSettings;
}

interface NetworkingService {
  id?: string;
  name?: string;
  type?: NetworkingServiceType;
  ipv4Config?: Ipv4Config;
  ipv6Config?: Ipv6Config;
  ipv6Privacy?: Ipv6Privacy;
  nameservers?: string[];
  searchDomains?: string[];
  timeservers?: string[];
  domain?: string;

  // From here everything only applies to wireless networks
  networkName?: string;
  ssid?: string;
  passphrase?: string;
  security?: Security;
  isHidden?: boolean;

  // From here everything only applies when security is Ieee8021x
  eap?: EapType;
  caCert?: string;
  caCertType?: CaCertType;
  privateKey?: string;
  privateKeyType?: PrivateKeyType;
  privateKeyPassphrase?: string;
  privateKeyPassphraseType?: PrivateKeyPassphraseType;
  identity?: string;
  anonymousIdentity?: string;
  subjectMatch?: string;
  altSubjectMatch?: string;
  domainSuffixMatch?: string;
  domainMatch?: string;
  phase2?: Phase2Type;
  isPhase2EapBased?: boolean;
}

interface NetworkingServicesPatch {
  services?: NetworkingService[];
}

const getNetworkingServices = () =>
  fetchAuthnd(`/networking/service`, {
    method: "GET",
  })
    .then((response) => response.json())
    .then((json) => json as Paginated<NetworkingService>);

const saveNetworkingServices = (patch: NetworkingServicesPatch) =>
  fetchAuthnd(`/user/${patch}`, {
    method: "PATCH",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(patch),
  });

interface NetworkSettingsProps extends SettingPageProps {}

const NetworkSettings: React.FC<NetworkSettingsProps> = (props) => {
  const { patch, setPatch } = useSetting<
    Paginated<NetworkingService>,
    NetworkingServicesPatch
  >({
    pageProps: props,
    isModified: (original, patch) => {
      return true;
    },
    fetch(): Promise<Paginated<NetworkingService>> {
      return getNetworkingServices();
    },
    initPatch(original: Paginated<NetworkingService>): NetworkingServicesPatch {
      return { services: original.data };
    },
    queryKey: "networkingServices",
    save: saveNetworkingServices,
    settingsKey: "networkingServices",
  });

  const updateService = (i: number, s: NetworkingService) => {
    const services = patch?.services || [];
    services[i] = s;
    setPatch({
      services: services,
    });
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        width: "100%",
        height: "100%",
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          width: "100%",
        }}
      >
        <div style={{ display: "flex", flexDirection: "column" }}>
          <Typography use="headline5">Netwerkinstellingen</Typography>
          <Typography use="headline6" style={{ marginTop: 8 }}>
            Netwerken
          </Typography>
        </div>

        <Button
          icon={"create"}
          style={{ backgroundColor: "#4CAF50" }}
          raised
          onClick={() => {
            setPatch({
              ...patch,
              services: [
                ...(patch?.services || []),
                {
                  name: "Nieuw netwerk",
                },
              ],
            });
          }}
        >
          Netwerk toevoegen
        </Button>
      </div>

      <List
        style={{
          borderTop:
            "1px solid var(--mdc-theme-text-hint-on-background, rgba(0, 0, 0, 0.38))",
          padding: 0,
          marginTop: 8,
          overflowY: "scroll",
          height: "100%",
        }}
      >
        {patch?.services?.map((service, i) => (
          <CollapsibleList
            style={{
              borderBottom:
                "1px solid var(--mdc-theme-text-hint-on-background, rgba(0, 0, 0, 0.38))",
            }}
            handle={
              <SimpleListItem
                text={service.name}
                metaIcon={<Icon icon="chevron_right" />}
              />
            }
          >
            <div
              style={{
                display: "flex",
                flexDirection: "column",
                margin: 16,
                marginTop: 0,
                height: "100%",
              }}
            >
              <TextField
                style={{
                  marginTop: 16,
                }}
                label={"Naam"}
                outlined
                value={service.name}
                onInput={(evt) => {
                  updateService(i, {
                    ...service,
                    name: (evt.target as HTMLInputElement).value,
                  });
                }}
              />

              <div
                style={{
                  marginTop: 16,
                }}
              >
                <Select
                  label={"Netwerktype"}
                  enhanced
                  outlined
                  options={{
                    [NetworkingServiceType.Ethernet]: "Ethernet",
                    [NetworkingServiceType.Wifi]: "Wi-Fi",
                  }}
                  onChange={(evt) => {
                    updateService(i, {
                      ...service,
                      type: (evt.target as HTMLSelectElement)
                        .value as NetworkingServiceType,
                    });
                  }}
                />
              </div>

              <div
                style={{
                  marginTop: 16,
                }}
              >
                <Select
                  label={"IPv4"}
                  enhanced
                  outlined
                  placeholder={"Automatisch"}
                  options={{
                    [Ipv4ConfigType.Off]: "Uit",
                    [Ipv4ConfigType.Dhcp]: "DHCP",
                    [Ipv4ConfigType.Custom]: "Handmatig",
                  }}
                  onChange={(evt) => {
                    updateService(i, {
                      ...service,
                      ipv4Config: {
                        ...service.ipv4Config,
                        type: (evt.target as HTMLSelectElement)
                          .value as Ipv4ConfigType,
                      },
                    });
                  }}
                />
              </div>

              {service?.ipv4Config?.type === Ipv4ConfigType.Custom && (
                <>
                  <TextField
                    style={{
                      marginTop: 16,
                    }}
                    label={"IPv4-netwerk"}
                    outlined
                  />

                  <TextField
                    style={{
                      marginTop: 16,
                    }}
                    label={"IPv4-netmask"}
                    outlined
                  />

                  <TextField
                    style={{
                      marginTop: 16,
                    }}
                    label={"IPv4-gateway"}
                    outlined
                  />
                </>
              )}

              <div
                style={{
                  marginTop: 16,
                }}
              >
                <Select
                  label={"IPv6"}
                  enhanced
                  outlined
                  placeholder={"Automatisch"}
                  options={{
                    [Ipv6ConfigType.Off]: "Uit",
                    [Ipv6ConfigType.Auto]: "Automatisch",
                    [Ipv6ConfigType.Custom]: "Handmatig",
                  }}
                  onChange={(evt) => {
                    updateService(i, {
                      ...service,
                      ipv6Config: {
                        ...service.ipv6Config,
                        type: (evt.target as HTMLSelectElement)
                          .value as Ipv6ConfigType,
                      },
                    });
                  }}
                />
              </div>

              {service?.ipv6Config?.type !== Ipv6ConfigType.Off && (
                <div
                  style={{
                    marginTop: 16,
                  }}
                >
                  <Select
                    label={"IPv6-privacy"}
                    enhanced
                    outlined
                    placeholder={"Automatisch"}
                    options={{
                      [Ipv6Privacy.Disabled]: "Uit",
                      [Ipv6Privacy.Enabled]: "Aan",
                      [Ipv6Privacy.Preferred]: "Bij voorkeur",
                    }}
                    onChange={(evt) => {
                      updateService(i, {
                        ...service,
                        ipv6Privacy: (evt.target as HTMLSelectElement)
                          .value as Ipv6Privacy,
                      });
                    }}
                  />
                </div>
              )}

              {service?.ipv6Config?.type === Ipv6ConfigType.Custom && (
                <>
                  <TextField
                    style={{
                      marginTop: 16,
                    }}
                    label={"IPv6-netwerk"}
                    outlined
                  />

                  <TextField
                    style={{
                      marginTop: 16,
                    }}
                    label={"IPv6-prefixlengte"}
                    type={"number"}
                    outlined
                  />

                  <TextField
                    style={{
                      marginTop: 16,
                    }}
                    label={"IPv6-gateway"}
                    outlined
                  />
                </>
              )}

              {service?.type === NetworkingServiceType.Wifi && (
                <>
                  <TextField
                    style={{
                      marginTop: 16,
                    }}
                    label={"Netwerknaam"}
                    outlined
                  />

                  <div
                    style={{
                      marginTop: 16,
                    }}
                  >
                    <Select
                      label={"Beveiliging"}
                      enhanced
                      outlined
                      options={{
                        [Security.Psk]: "WPA",
                        [Security.Ieee8021x]: "WPA-Enterprise",
                        [Security.Wep]: "WEP",
                      }}
                      onChange={(evt) => {
                        updateService(i, {
                          ...service,
                          security: (evt.target as HTMLSelectElement)
                            .value as Security,
                        });
                      }}
                    />
                  </div>

                  <TextField
                    style={{
                      marginTop: 16,
                    }}
                    label={"Wachtwoord"}
                    outlined
                  />

                  <Switch
                    style={{
                      marginTop: 16,
                    }}
                  >
                    Netwerk is verborgen
                  </Switch>

                  {service?.security === Security.Ieee8021x && (
                    <>
                      <TextField
                        style={{
                          marginTop: 16,
                        }}
                        label={"Identiteit"}
                        outlined
                      />

                      <TextField
                        style={{
                          marginTop: 16,
                        }}
                        label={"Anonieme identiteit"}
                        outlined
                      />
                    </>
                  )}
                </>
              )}

              <Button danger raised icon={"delete"} style={{ marginTop: 16 }}>
                Netwerk verwijderen
              </Button>
            </div>
          </CollapsibleList>
        ))}
      </List>
    </div>
  );
};

export default NetworkSettings;
