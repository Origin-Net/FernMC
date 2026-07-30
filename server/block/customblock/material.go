package customblock


type Material struct {
	
	texture string
	
	renderMethod Method
	
	faceDimming bool
	
	ambientOcclusion bool
}



func NewMaterial(texture string, method Method) Material {
	return Material{
		texture:          texture,
		renderMethod:     method,
		faceDimming:      true,
		ambientOcclusion: method.AmbientOcclusion(),
	}
}


func (m Material) WithFaceDimming() Material {
	m.faceDimming = true
	return m
}


func (m Material) WithoutFaceDimming() Material {
	m.faceDimming = false
	return m
}


func (m Material) WithAmbientOcclusion() Material {
	m.ambientOcclusion = true
	return m
}


func (m Material) WithoutAmbientOcclusion() Material {
	m.ambientOcclusion = false
	return m
}


func (m Material) Encode() map[string]any {
	return map[string]any{
		"texture":           m.texture,
		"render_method":     m.renderMethod.String(),
		"face_dimming":      m.faceDimming,
		"ambient_occlusion": m.ambientOcclusion,
	}
}
