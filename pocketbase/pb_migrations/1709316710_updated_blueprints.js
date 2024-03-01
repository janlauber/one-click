/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("vs5gr49hpah1g9q")

  collection.listRule = "@request.auth.id != \"\" && (private = false && @request.auth.id = owner.id)"
  collection.viewRule = "@request.auth.id != \"\" && (private = false && @request.auth.id = owner.id)"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("vs5gr49hpah1g9q")

  collection.listRule = "@request.auth.id != \"\" && (private = true)"
  collection.viewRule = "@request.auth.id != \"\""

  return dao.saveCollection(collection)
})
