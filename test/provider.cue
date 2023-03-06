package test

#Apply: {
	#do:       "apply"
	#provider: "test"
	$params: {
		cluster: string
		resource: {
			...
		}
		options: {
			threeWayMergePatch: {
				enabled:          bool
				annotationPrefix: string
			}
		}
	}
	$returns: {
		...
	}
}

#Get: {
	#do:       "get"
	#provider: "test"
	$params: {
		cluster: string
		resource: {
			...
		}
		options: {
			threeWayMergePatch: {
				enabled:          bool
				annotationPrefix: string
			}
		}
	}
	$returns: {
		...
	}
}

#List: {
	#do:       "list"
	#provider: "test"
	$params: {
		cluster: string
		filter?: null | {
			namespace?: string
			matchingLabels?: {
				[string]: string
			}
		}
		resource: {
			...
		}
	}
	$returns: {
		...
	}
}

#Patch: {
	#do:       "patch"
	#provider: "test"
	$params: {
		cluster: string
		resource: {
			...
		}
		patch: {
			type: string
			data: {
				...
			}
		}
	}
	$returns: {
		...
	}
}
