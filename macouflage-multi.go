package main

import (
	"strings"
	"fmt"
	"log"
	lmf "github.com/subgraph/libmacouflage"
)


func getCurrentMacInfo(name string) (result string, err error) {
	currentMacInfo, err := getMacInfo(name, "CurrentMAC")
	if err != nil {
		return
	}
	permanentMac, err := lmf.GetPermanentMac(name)
	if err != nil {
		LogWriter.Err(err.Error())
	}
	permanentMacVendor, err := lmf.FindVendorByMac(permanentMac.String())
	if err != nil {
		if strings.HasPrefix(err.Error(),
			"No vendor found in OuiDb for vendor prefix") {
			permanentMacVendor.Vendor = "Unknown"
		} else {
			return
		}
	}
	if privatevar {
		result = fmt.Sprintf(
			"%s - CurrentMAC: xx:xx:xx:xx:xx:xx (Unknown) Permanent MAC: xx:xx:xx:xx:xx:xx (Unknown)",
			name)
	} else {
		result = fmt.Sprintf("%sPermanent MAC: %s (%s)",
			currentMacInfo, permanentMac, permanentMacVendor.Vendor)
	}
	return
}

func getMacInfo(name string, macType string) (result string, err error) {
	newMac, err := lmf.GetCurrentMac(name)
	if err != nil {
		return
	}
	newMacVendor, err := lmf.FindVendorByMac(newMac.String())
	if err != nil {
		if err == err.(lmf.NoVendorError) {
			newMacVendor.Vendor = "Unknown"
			err = nil
		} else {
			return
		}
	}
	result = fmt.Sprintf("%s - %s: %s (%s) ",
		name, macType, newMac, newMacVendor.Vendor)
	return
}

func listVendors(keyword string, isPopular bool) (results string, err error) {
	var ouis []lmf.Oui
	var vendors []string

	if isPopular {
		ouis, err = lmf.FindAllPopularOuis()
		if err != nil {
			return
		}
	} else {
		ouis, err = lmf.FindVendorsByKeyword(keyword)
		if err != nil {
			return
		}
	}
	if len(ouis) == 0 {
		results = fmt.Sprintf("No vendors found in search.")
		return
	} else {
		vendors = append(vendors, fmt.Sprintf("#\tVendorPrefix\tVendor"))
		for i, result := range ouis {
			vendors = append(vendors, fmt.Sprintf("%d\t%s\t%s", i+1,
				result.VendorPrefix, result.Vendor))
		}
		results = strings.Join(vendors, "\n")
	}
	return
}

func spoofMacEnding(name string) (err error) {
	currentMacInfo, err := getCurrentMacInfo(name)
	if err != nil {
		return
	}
	LogWriter.Info(currentMacInfo)
	changed, err := lmf.SpoofMacSameVendor(name, true)
	if err != nil {
		return
	}
	if changed {
		newMac, err2 := getMacInfo(name, "New MAC")
		if err2 != nil {
			err = err2
			return
		}
		LogWriter.Info(newMac)
	}
	return
}

func spoofMacAnother(name string) (err error) {
	currentMacInfo, err := getCurrentMacInfo(name)
	if err != nil {
		return
	}
	LogWriter.Info(currentMacInfo)
	changed, err := lmf.SpoofMacSameDeviceType(name)
	if err != nil {
		return
	}
	if changed {
		newMac, err2 := getMacInfo(name, "New MAC")
		if err2 != nil {
			err = err2
			return
		}
		LogWriter.Info(newMac)
	}
	return
}

func spoofMacAny(name string) (err error) {
	currentMacInfo, err := getCurrentMacInfo(name)
	if err != nil {
		return
	}
	LogWriter.Info(currentMacInfo)
	changed, err := lmf.SpoofMacAnyDeviceType(name)
	if err != nil {
		return
	}
	if changed {
		newMac, err2 := getMacInfo(name, "New MAC")
		if err2 != nil {
			err = err2
			return
		}
		LogWriter.Info(newMac)
	}
	return
}

func revertMac(name string) (err error) {
	currentMacInfo, err := getCurrentMacInfo(name)
	if err != nil {
		return
	}
	LogWriter.Info(currentMacInfo)
	permMac, err := lmf.GetPermanentMac(name)
	if err != nil {
		return
	}
	err = lmf.RevertMac(name)
	if err != nil {
		return
	}
	newMac, err := lmf.GetCurrentMac(name)
	if err != nil {
		return
	}
	if lmf.CompareMacs(permMac, newMac) {
		newMac, err2 := getMacInfo(name, "New MAC")
		if err2 != nil {
			err = err2
			return
		}
		LogWriter.Info(newMac)
	}
	return
}

func spoofMacRandom(name string, bia bool) (err error) {
	currentMacInfo, err := getCurrentMacInfo(name)
	if err != nil {
		return
	}
	LogWriter.Info(currentMacInfo)
	changed, err := lmf.SpoofMacRandom(name, bia)
	if err != nil {
		return
	}
	if changed {
		newMac, err2 := getMacInfo(name, "New MAC")
		if err2 != nil {
			err = err2
			return
		}
		LogWriter.Info(newMac)
	}
	return
}

func spoofMacPopular(name string) (err error) {
	currentMacInfo, err := getCurrentMacInfo(name)
	if err != nil {
		return
	}
	LogWriter.Info(currentMacInfo)
	changed, err := lmf.SpoofMacPopular(name)
	if err != nil {
		return
	}
	if changed {
		newMac, err2 := getMacInfo(name, "New MAC")
		if err2 != nil {
			err = err2
			return
		}
		LogWriter.Info(newMac)
	}
	return
}

func spoofAll(mode string) (spoofErrors []error) {
	ifaces, err := lmf.GetInterfaces()
	if err != nil {
		log.Fatalf("Could not get list of interfaces: %s", err)
	}
	for _, iface := range ifaces {
		changed, err := lmf.MacChanged(iface.Name)
		if !changed || forcevar {
			switch mode {
			case "ending":
				err = spoofMacEnding(iface.Name)
			case "another":
				err = spoofMacAnother(iface.Name)
			case "any":
				err = spoofMacAny(iface.Name)
			case "random":
				err = spoofMacRandom(iface.Name, true)
			case "popular":
				err = spoofMacPopular(iface.Name)
			default:
				err = fmt.Errorf("Invalid mode: %s", mode)
			}
			if err != nil {
				spoofErrors = append(spoofErrors, err)
			}
		}
	}
	return
}