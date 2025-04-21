package domain_repos

import (
	"github.com/SabianF/ghasp/src/common/presentation/components"
)

func NewLayoutPropsDefault() components.LayoutProps {
	return components.LayoutProps{
		IncludeSidebar: true,
		SidebarMenuLayoutProps: components.SidebarMenuLayoutProps{
			MenuProps: components.SidebarMenuProps{
				SidebarMenuItemProps: []components.SidebarMenuItemProps{
					{ Text: RootPageName, Url: RootPageUrl },
					{ Text: HtmxExamplesPageName, Url: HtmxExamplesPageUrl },
				},
			},
		},
	}
}

func NewLayoutPropsNoSidebar() components.LayoutProps {
	return components.LayoutProps{IncludeSidebar: false}
}
