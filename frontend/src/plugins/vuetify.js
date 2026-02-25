// Vuetify 配置文件
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import '@mdi/font/css/materialdesignicons.css'

// 前台深色主题（金色点缀）
const frontendDarkTheme = {
  dark: true,
  colors: {
    background: '#0a0a0a',
    surface: '#1a1a1a',
    primary: '#c9a962',
    'primary-darken-1': '#b08d4f',
    'primary-lighten-1': '#d4b978',
    secondary: '#8a8a8a',
    'secondary-darken-1': '#6a6a6a',
    accent: '#d4a574',
    error: '#cf6679',
    info: '#2196f3',
    success: '#4caf50',
    warning: '#fb8c00',
  },
  variables: {
    'border-opacity': 0.12,
    'high-emphasis-opacity': 0.87,
    'medium-emphasis-opacity': 0.60,
    'disabled-opacity': 0.38,
    'idle-opacity': 0.10,
    'activated-opacity': 0.12,
  }
}

// 后台浅色主题（Material Design）
const adminLightTheme = {
  dark: false,
  colors: {
    background: '#FFFFFF',
    surface: '#FFFFFF',
    primary: '#1976D2',
    'primary-darken-1': '#1565C0',
    'primary-lighten-1': '#1E88E5',
    secondary: '#757575',
    'secondary-darken-1': '#616161',
    accent: '#82B1FF',
    error: '#FF5252',
    info: '#2196F3',
    success: '#4CAF50',
    warning: '#FF9800',
  },
  variables: {
    'border-opacity': 0.12,
    'high-emphasis-opacity': 0.87,
    'medium-emphasis-opacity': 0.60,
    'disabled-opacity': 0.38,
    'idle-opacity': 0.04,
    'activated-opacity': 0.12,
  }
}

// 创建 Vuetify 实例
export default createVuetify({
  components,
  directives,

  // 主题配置
  theme: {
    defaultTheme: 'frontendDarkTheme',
    themes: {
      frontendDarkTheme,
      adminLightTheme,
    },
  },

  // 默认组件配置
  defaults: {
    VBtn: {
      // 保持按钮文字不大写（前台风格）
      style: 'text-transform: none;',
    },
    VCard: {
      // 卡片圆角
      rounded: 'lg',
    },
    VTextField: {
      // 输入框样式
      variant: 'outlined',
      density: 'comfortable',
    },
    VTextarea: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VCombobox: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VFileInput: {
      variant: 'outlined',
      density: 'comfortable',
    },
  },

  // 图标配置
  icons: {
    defaultSet: 'mdi',
  },
})
