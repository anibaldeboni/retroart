package theme

// Theme é a interface principal para acessar o design system
type Theme interface {
	GetDesignSystem() DesignSystem
	GetButtonStyle(styleType ComponentStyleType) ButtonStyle
	GetCheckboxListStyle() CheckboxListStyle
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
