docker_build('rgeraskin/dogoncall',
    context='.',
)
k8s_yaml(helm(
    './helm',
    name='dogoncall',
    values=['./helm/values-private.yaml'],
))
