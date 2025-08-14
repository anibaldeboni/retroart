package theme

// Theme é a interface principal para acessar o design system
type Theme interface {
	GetDesignSystem() DesignSystem
	GetButtonStyle(styleType ComponentStyleType) ButtonStyle
	GetCheckboxListStyle() CheckboxListStyle
	GetInputTextStyle() InputTextStyle
	GetVirtualKeyboardStyle() VirtualKeyboardStyle
	GetMainContainerStyle() ContainerStyle
	GetContentContainerStyle() ContainerStyle
}

// DefaultTheme implementa a interface Theme com o design system padrão
type DefaultTheme struct {
	designSystem DesignSystem
}

// NewDefaultTheme cria uma nova instância do tema padrão
func NewDefaultTheme() Theme {
	return &DefaultTheme{
		designSystem: DefaultDesignSystem(),
	}
}

// GetDesignSystem retorna o design system completo
func (t *DefaultTheme) GetDesignSystem() DesignSystem {
	return t.designSystem
}

// GetButtonStyle retorna o estilo para botões
func (t *DefaultTheme) GetButtonStyle(styleType ComponentStyleType) ButtonStyle {
	return t.designSystem.GetButtonStyle(styleType)
}

// GetCheckboxListStyle retorna o estilo para checkbox lists
func (t *DefaultTheme) GetCheckboxListStyle() CheckboxListStyle {
	return t.designSystem.GetCheckboxListStyle()
}

// GetInputTextStyle retorna o estilo para campos de texto
func (t *DefaultTheme) GetInputTextStyle() InputTextStyle {
	return t.designSystem.GetInputTextStyle()
}

// GetVirtualKeyboardStyle retorna o estilo para o teclado virtual
func (t *DefaultTheme) GetVirtualKeyboardStyle() VirtualKeyboardStyle {
	return t.designSystem.GetVirtualKeyboardStyle()
}

// GetMainContainerStyle retorna o estilo para container principal
func (t *DefaultTheme) GetMainContainerStyle() ContainerStyle {
	return t.designSystem.GetMainContainerStyle()
}

// GetContentContainerStyle retorna o estilo para containers de conteúdo
func (t *DefaultTheme) GetContentContainerStyle() ContainerStyle {
	return t.designSystem.GetContentContainerStyle()
}

// Instância global do tema (singleton)
var currentTheme Theme

// GetCurrentTheme retorna o tema atual
func GetCurrentTheme() Theme {
	if currentTheme == nil {
		currentTheme = NewDefaultTheme()
	}
	return currentTheme
}

// SetTheme define um novo tema
func SetTheme(theme Theme) {
	currentTheme = theme
}

// Funções de conveniência para acessar rapidamente o tema atual

// GetButtonStyle é uma função de conveniência para obter estilos de botão
func GetButtonStyle(styleType ComponentStyleType) ButtonStyle {
	return GetCurrentTheme().GetButtonStyle(styleType)
}

// GetCheckboxListStyle é uma função de conveniência para obter estilos de checkbox list
func GetCheckboxListStyle() CheckboxListStyle {
	return GetCurrentTheme().GetCheckboxListStyle()
}

// GetInputTextStyle é uma função de conveniência para obter estilos de campo de texto
func GetInputTextStyle() InputTextStyle {
	return GetCurrentTheme().GetInputTextStyle()
}

// GetVirtualKeyboardStyle é uma função de conveniência para obter estilos de teclado virtual
func GetVirtualKeyboardStyle() VirtualKeyboardStyle {
	return GetCurrentTheme().GetVirtualKeyboardStyle()
}

// GetMainContainerStyle é uma função de conveniência para obter estilos de container principal
func GetMainContainerStyle() ContainerStyle {
	return GetCurrentTheme().GetMainContainerStyle()
}

// GetContentContainerStyle é uma função de conveniência para obter estilos de container de conteúdo
func GetContentContainerStyle() ContainerStyle {
	return GetCurrentTheme().GetContentContainerStyle()
}

// GetColors é uma função de conveniência para obter a paleta de cores
func GetColors() ColorPalette {
	return GetCurrentTheme().GetDesignSystem().Colors
}

// GetTypography é uma função de conveniência para obter configurações de tipografia
func GetTypography() Typography {
	return GetCurrentTheme().GetDesignSystem().Typography
}

// GetSpacing é uma função de conveniência para obter configurações de espaçamento
func GetSpacing() Spacing {
	return GetCurrentTheme().GetDesignSystem().Spacing
}
