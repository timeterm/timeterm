import QtQuick 2.0
import QtQml.Models 2.3

import Timeterm.Api.Zermelo 1.0

ZermeloApopintments {
    data: [
        ZermeloApopintment {
            subjects: ["nat", "nat1"]
            groups: ["gv5.nat1", "gv6.nat3"]
            locations: ["g208", "g209"]
            teachers: ["mrd", "abc"]
        },
        ZermeloApopintment {
            subjects: ["nat", "nat1"]
            groups: ["gv5.nat1", "gv6.nat3"]
            locations: ["g208", "g209"]
            teachers: ["mrd", "abc"]
        },
        ZermeloApopintment {
            subjects: ["nat", "nat1"]
            groups: ["gv5.nat1", "gv6.nat3"]
            locations: ["g208", "g209"]
            teachers: ["mrd", "abc"]
        },
        ZermeloApopintment {
            subjects: ["nat", "nat1"]
            groups: ["gv5.nat1", "gv6.nat3"]
            locations: ["g208", "g209"]
            teachers: ["mrd", "abc"]
        }
    ]

    //ListModel {
    //    ListElement {
    //        subjects: ["nat", "nat2"]
    //        teacher: "MRD"
    //        group: "gv5.nat3"
    //        location: "g208"
    //    }
    //    ListElement {
    //        subjects: [
    //            ListElement {
    //                subject: "nat"
    //            }
    //        ]
    //        teacher: "MRD"
    //        group: "gv5.nat3"
    //        location: "g208"
    //    }
    //    ListElement {
    //        subjects: [
    //            ListElement {
    //                subject: "nat"
    //            }
    //        ]
    //        teacher: "MRD"
    //        group: "gv5.nat3"
    //        location: "g208"
    //    }
    //    ListElement {
    //        subjects: [
    //            ListElement {
    //                subject: "nat"
    //            }
    //        ]
    //        teacher: "MRD"
    //        group: "gv5.nat3"
    //        location: "g208"
    //    }
    //}
}
