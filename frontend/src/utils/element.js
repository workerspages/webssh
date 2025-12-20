// 导入自己需要的组件
import {
    Form,
    FormItem,
    Dialog,
    Row,
    Col,
    Button,
    ButtonGroup,
    Table,
    TableColumn,
    Input,
    Message,
    Container,
    Header,
    Main,
    Upload,
    Dropdown,
    DropdownMenu,
    DropdownItem,
    Tabs,
    TabPane,
    Divider,
    Tooltip,
    Menu,
    Submenu,
    MenuItem,
    MenuItemGroup,
    Aside,
    Breadcrumb,
    BreadcrumbItem
} from 'element-ui'
const element = {
    install: function (Vue) {
        Vue.use(Input)
        Vue.use(Dialog)
        Vue.use(Row)
        Vue.use(Col)
        Vue.use(Form)
        Vue.use(FormItem)
        Vue.use(Button)
        Vue.use(ButtonGroup)
        Vue.use(Table)
        Vue.use(TableColumn)
        Vue.use(Container)
        Vue.use(Header)
        Vue.use(Aside)
        Vue.use(Main)
        Vue.use(Menu)
        Vue.use(Submenu)
        Vue.use(MenuItem)
        Vue.use(MenuItemGroup)
        Vue.use(Upload)
        Vue.use(Dropdown)
        Vue.use(DropdownItem)
        Vue.use(DropdownMenu)
        Vue.use(Tabs)
        Vue.use(TabPane)
        Vue.use(Divider)
        Vue.use(Tooltip)
        Vue.use(Breadcrumb)
        Vue.use(BreadcrumbItem)
        Vue.prototype.$message = Message
    }
}
export default element