/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("w9y791qv1ymfu3m")

  collection.viewRule = "@request.auth.id != \"\" && (@request.auth.id = blueprint.owner.id || @request.auth.id ?= blueprint.users.id || @request.auth.id = project.user.id)"
  collection.createRule = "@request.auth.id != \"\" && (@request.auth.id = blueprint.owner.id || @request.auth.id = project.user.id)"
  collection.updateRule = "@request.auth.id != \"\" && (@request.auth.id = blueprint.owner.id || @request.auth.id = project.user.id)"
  collection.deleteRule = "@request.auth.id != \"\" && (@request.auth.id = blueprint.owner.id || @request.auth.id = project.user.id)"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("w9y791qv1ymfu3m")

  collection.viewRule = ""
  collection.createRule = null
  collection.updateRule = null
  collection.deleteRule = null

  return dao.saveCollection(collection)
})
