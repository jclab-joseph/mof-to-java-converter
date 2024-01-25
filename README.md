# mof2java converter

## Example

### Generating

**mapper.txt :**

```text
wsdl="http://schemas.xmlsoap.org/wsdl/"
xsd="http://www.w3.org/2001/XMLSchema"
wsman="http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd"
soap="http://www.w3.org/2003/05/soap-envelope"
wscat="http://schemas.xmlsoap.org/ws/2005/06/wsmancat"
wsa="http://schemas.xmlsoap.org/ws/2004/08/addressing"
wxf="http://schemas.xmlsoap.org/ws/2004/09/transfer"
wse="http://schemas.xmlsoap.org/ws/2004/08/eventing"
cim="http://schemas.dmtf.org/wbem/wscim/1/common"
wsen="http://schemas.xmlsoap.org/ws/2004/09/enumeration"
vxiny="http://visualcim.org/1/common/Common_Types"
vxgrl="http://visualcim.org/1/common/WsmanGeneral"
idnty="http://schemas.dmtf.org/wbem/wsman/identity/1/wsmanidentity.xsd"
CIM_BIOSFeature="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_BIOSFeature"
CIM_IEEE8021xCapabilities="http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_IEEE8021xCapabilities"
```

**mof-list.txt :**

```text
CIM_BIOSFeature.mof
CIM_IEEE8021xCapabilities.mof
```

**Run:**

```bash
mkdir -p $PWD/output/
mof2java-converter.exe \
  --package "org.example.something.mof" \
  --ns-mapper-file "mapper.txt" \
  --paths /somewhere/omi/Unix/share/omischema/CIM-2.32.0:$PWD/mofs/ \
  --srcdir $PWD/mofs/ \
  --outdir $PWD/output/ \
  $(cat mof-list.txt)
```

### Generated

**package-info.java :**

```java
@XmlSchema(
        elementFormDefault = XmlNsForm.QUALIFIED
)
package org.example.something.mof;

import jakarta.xml.bind.annotation.XmlNsForm;
import jakarta.xml.bind.annotation.XmlSchema;
```

**CIM_BIOSFeature.java :**

```java
package org.example.something.mof;

import jakarta.xml.bind.annotation.XmlAccessType;
import jakarta.xml.bind.annotation.XmlAccessorType;
import jakarta.xml.bind.annotation.XmlElement;
import jakarta.xml.bind.annotation.XmlType;
import jakarta.xml.bind.annotation.adapters.XmlJavaTypeAdapter;
import javax.cim.UnsignedInteger16;

@XmlAccessorType(XmlAccessType.PROPERTY)
@XmlType(namespace = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_BIOSFeature")
public class CIM_BIOSFeature extends CIM_SoftwareFeature {
    public static final String RESOURCE_URI = "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_BIOSFeature";
    
    /**
     * An array of integers that specify the features supported by the BIOS. For example, one can specify that PnP capabilities are provided (value=9) or that infrared devices are supported (21). Values specified in the enumeration are taken from both DMI and SMBIOS (the Type 0 structure, the BIOS Characteristics and BIOS Characteristics Extension Bytes attributes.
     * 
    * ValuesMap: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 160]
    * Values: [Other, Unknown, Undefined, ISA Support, MCA Support, EISA Support, PCI Support, PCMCIA Support, PnP Support, APM Support, Upgradeable BIOS, BIOS Shadowing Allowed, VL VESA Support, ESCD Support, LS-120 Boot Support, ACPI Support, I2O Boot Support, USB Legacy Support, AGP Support, PC Card, IR, 1394, I2C, Smart Battery, ATAPI ZIP Drive Boot Support, 1394 Boot Support, Boot from CD, Selectable Boot, BIOS ROM is Socketed, Boot from PCMCIA, EDD Specification Support, PC-98]
     */
    @XmlElement(name = "Characteristics")
    @XmlJavaTypeAdapter(UnsignedInteger16.XmlAdapter.class)
    public java.util.List<UnsignedInteger16> Characteristics;
    
    
    /**
     * An array of free-form strings providing more detailed explanations for any of the BIOS features indicated in the Characteristics array. Note, each entry of this array is related to the entry in the Characteristics array that is located at the same index.
     */
    @XmlElement(name = "CharacteristicDescriptions")
    public java.util.List<String> CharacteristicDescriptions;
}
```

UnsignedInteger16, UnsignedInteger32, etc are at https://github.com/jc-lab/abstract-wsman.

## LICENSE

[GPL-3.0](./LICENSE)
