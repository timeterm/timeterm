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
import deepEqual from "deep-equal";
import { v4 as uuidv4 } from "uuid";
import useFileImporter from "./useFileImporter";

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
  Ieee8021x = "Ieee8021x",
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
  deleted?: string[];
}

const getNetworkingServices: (
  offset?: number
) => Promise<NetworkingService[]> = async (offset: number = 0) => {
  const rsp = await fetchAuthnd(`/networking/service?offset=${offset}`, {
    method: "GET",
  })
    .then((response) => response.json())
    .then((json) => json as Paginated<NetworkingService>);

  if (rsp.offset + rsp.maxAmount < rsp.total) {
    return [
      ...rsp.data,
      ...(await getNetworkingServices(rsp.offset + rsp.maxAmount)),
    ];
  }

  return rsp.data;
};

const saveNetworkingServices = async (patch: NetworkingServicesPatch) => {
  for (const service of patch.services || []) {
    if (service.id?.startsWith("new-")) {
      const uuidlessService = {
        ...service,
        id: null,
      };
      await fetchAuthnd(`/networking/service`, {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(uuidlessService),
      });
    } else {
      await fetchAuthnd(`/networking/service/${service.id}`, {
        method: "PUT",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(service),
      });
    }
  }

  for (const id of patch.deleted || []) {
    if (!id.startsWith("new-")) {
      await fetchAuthnd(`/networking/service/${id}`, {
        method: "DELETE",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      });
    }
  }
};

interface NetworkSettingsProps extends SettingPageProps {}

const NetworkSettings: React.FC<NetworkSettingsProps> = (props) => {
  const { patch, setPatch } = useSetting<
    NetworkingService[],
    NetworkingServicesPatch
  >({
    pageProps: props,
    isModified: (original, patch) => {
      return !deepEqual(original, patch.services);
    },
    fetch(): Promise<NetworkingService[]> {
      return getNetworkingServices();
    },
    initPatch(original: NetworkingService[]): NetworkingServicesPatch {
      return { services: [...original], deleted: [] };
    },
    queryKey: "networkingServices",
    save: saveNetworkingServices,
    settingsKey: "networkingServices",
  });

  const updateService = (i: number, s: NetworkingService) => {
    const services = patch?.services || [];
    services[i] = s;
    setPatch({
      ...patch,
      services: services,
    });
  };

  const deleteService = (i: number) => {
    const services = patch?.services || [];
    if (services.length <= i) {
      return;
    }
    const service = patch?.services?.[i];

    setPatch({
      ...patch,
      services: [
        ...services.slice(0, i),
        ...services.slice(i + 1, services.length),
      ],
      deleted: [
        ...(patch?.deleted || []),
        ...(service?.id ? [service.id] : []),
      ],
    });
  };

  const caCertImporter = useFileImporter<number>(
    [".pem", ".der"],
    (filename, contents, svci) => {
      if (svci !== undefined && (patch?.services?.length || 0) > svci) {
        console.log("TTT");
        const type = filename.endsWith(".pem")
          ? CaCertType.Pem
          : filename.endsWith(".der")
          ? CaCertType.Der
          : undefined;

        console.log("converting");
        const b64Cert = btoa(
          new Uint8Array(contents as ArrayBuffer).reduce(
            (data, byte) => data + String.fromCharCode(byte),
            ""
          )
        );
        updateService(svci, {
          ...patch?.services?.[svci],
          caCert: b64Cert,
          caCertType: type,
        });
      }
    }
  );

  const privateKeyImporter = useFileImporter<number>(
    [".pfx", ".pem", ".der"],
    (filename, contents, svci) => {
      if (svci !== undefined && (patch?.services?.length || 0) > svci) {
        const type = filename.endsWith(".pfx")
          ? PrivateKeyType.Pfx
          : filename.endsWith(".pem")
          ? PrivateKeyType.Pem
          : filename.endsWith(".der")
          ? PrivateKeyType.Der
          : undefined;

        const b64Key = btoa(
          new Uint8Array(contents as ArrayBuffer).reduce(
            (data, byte) => data + String.fromCharCode(byte),
            ""
          )
        );

        updateService(svci, {
          ...patch?.services?.[svci],
          privateKey: b64Key,
          privateKeyType: type,
        });
      }
    }
  );

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
                  // Generate a UUID (even though the backend assigns another one)
                  // Doing this prevents the service below this one from opening after the deletion of this one
                  // (when using the index as the key in the list instead of the id)
                  id: "new-" + uuidv4(),
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
        {patch?.services?.map((service, i) => {
          return (
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
              key={service.id}
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
                      onChange={(evt) =>
                        updateService(i, {
                          ...service,
                          ipv4Config: {
                            ...service.ipv4Config,
                            settings: {
                              ...service.ipv4Config?.settings,
                              network: (evt.target as HTMLInputElement).value,
                            },
                          },
                        })
                      }
                    />

                    <TextField
                      style={{
                        marginTop: 16,
                      }}
                      label={"IPv4-netmask"}
                      outlined
                      onChange={(evt) =>
                        updateService(i, {
                          ...service,
                          ipv4Config: {
                            ...service.ipv4Config,
                            settings: {
                              ...service.ipv4Config?.settings,
                              netmask: (evt.target as HTMLInputElement).value,
                            },
                          },
                        })
                      }
                    />

                    <TextField
                      style={{
                        marginTop: 16,
                      }}
                      label={"IPv4-gateway"}
                      outlined
                      onChange={(evt) =>
                        updateService(i, {
                          ...service,
                          ipv4Config: {
                            ...service.ipv4Config,
                            settings: {
                              ...service.ipv4Config?.settings,
                              gateway: (evt.target as HTMLInputElement).value,
                            },
                          },
                        })
                      }
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
                      onChange={(evt) =>
                        updateService(i, {
                          ...service,
                          ipv6Config: {
                            ...service.ipv6Config,
                            settings: {
                              ...service.ipv6Config?.settings,
                              network: (evt.target as HTMLInputElement).value,
                            },
                          },
                        })
                      }
                    />

                    <TextField
                      style={{
                        marginTop: 16,
                      }}
                      label={"IPv6-prefixlengte"}
                      type={"number"}
                      outlined
                      onChange={(evt) =>
                        updateService(i, {
                          ...service,
                          ipv6Config: {
                            ...service.ipv6Config,
                            settings: {
                              ...service.ipv6Config?.settings,
                              prefixLength: Number(
                                (evt.target as HTMLInputElement).value
                              ),
                            },
                          },
                        })
                      }
                    />

                    <TextField
                      style={{
                        marginTop: 16,
                      }}
                      label={"IPv6-gateway"}
                      outlined
                      onChange={(evt) =>
                        updateService(i, {
                          ...service,
                          ipv6Config: {
                            ...service.ipv6Config,
                            settings: {
                              ...service.ipv6Config?.settings,
                              gateway: (evt.target as HTMLInputElement).value,
                            },
                          },
                        })
                      }
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
                      onChange={(evt) =>
                        updateService(i, {
                          ...service,
                          networkName: (evt.target as HTMLInputElement).value,
                        })
                      }
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
                      onChange={(evt) =>
                        updateService(i, {
                          ...service,
                          passphrase: (evt.target as HTMLInputElement).value,
                        })
                      }
                    />

                    <Switch
                      style={{
                        marginTop: 16,
                      }}
                      onChange={(evt) =>
                        updateService(i, {
                          ...service,
                          isHidden: (evt.target as HTMLInputElement).checked,
                        })
                      }
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
                          onChange={(evt) =>
                            updateService(i, {
                              ...service,
                              identity: (evt.target as HTMLInputElement).value,
                            })
                          }
                        />

                        <TextField
                          style={{
                            marginTop: 16,
                          }}
                          label={"Anonieme identiteit"}
                          outlined
                          onChange={(evt) =>
                            updateService(i, {
                              ...service,
                              anonymousIdentity: (evt.target as HTMLInputElement)
                                .value,
                            })
                          }
                        />

                        <div
                          style={{
                            marginTop: 16,
                          }}
                        >
                          <Select
                            label={"EAP"}
                            enhanced
                            outlined
                            options={{
                              [EapType.Peap]: "PEAP",
                              [EapType.Tls]: "TLS",
                              [EapType.Ttls]: "TTLS",
                            }}
                            onChange={(evt) => {
                              updateService(i, {
                                ...service,
                                eap: (evt.target as HTMLSelectElement)
                                  .value as EapType,
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
                            label={"Phase 2"}
                            enhanced
                            outlined
                            options={{
                              [Phase2Type.Gtc]: "GTC",
                              [Phase2Type.MschapV2]: "MSCHAPv2",
                            }}
                            onChange={(evt) => {
                              updateService(i, {
                                ...service,
                                phase2: (evt.target as HTMLSelectElement)
                                  .value as Phase2Type,
                              });
                            }}
                          />
                        </div>

                        {service?.eap === EapType.Ttls && (
                          <>
                            <Switch
                              style={{
                                marginTop: 16,
                              }}
                              onChange={(evt) =>
                                updateService(i, {
                                  ...service,
                                  isPhase2EapBased: (evt.target as HTMLInputElement)
                                    .checked,
                                })
                              }
                            >
                              Phase 2 maakt gebruik van EAP
                            </Switch>
                          </>
                        )}

                        <Button
                          style={{
                            marginTop: 16,
                          }}
                          raised
                          onClick={() => caCertImporter(i)}
                        >
                          CA-certificaat uploaden
                        </Button>

                        <Button
                          style={{
                            marginTop: 16,
                          }}
                          raised
                          onClick={() => privateKeyImporter(i)}
                        >
                          Privésleutel uploaden
                        </Button>
                      </>
                    )}
                  </>
                )}

                <Button
                  danger
                  raised
                  icon={"delete"}
                  style={{ marginTop: 16 }}
                  onClick={() => deleteService(i)}
                >
                  Netwerk verwijderen
                </Button>
              </div>
            </CollapsibleList>
          );
        })}
      </List>
    </div>
  );
};

export default NetworkSettings;
